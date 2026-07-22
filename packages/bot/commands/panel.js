import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  ChannelType,
  EmbedBuilder,
  MessageFlags,
  ModalBuilder,
  PermissionFlagsBits,
  SlashCommandBuilder,
  TextInputBuilder,
  TextInputStyle,
} from "discord.js";
import { env } from "../config/env.js";
import {
  addAudit,
  auditLog,
  getPanelSettings,
  incidents,
  resetPanelSettings,
  tickets,
  updatePanelSettings,
  warnings,
} from "../utils/store.js";

const PANEL_REFRESH_BUTTON_ID = "panel:refresh";
const PANEL_TICKETS_BUTTON_ID = "panel:tickets";
const PANEL_INCIDENTS_BUTTON_ID = "panel:incidents";
const PANEL_REPORT_BUTTON_ID = "panel:report";
const PANEL_ANNOUNCE_BUTTON_ID = "panel:announce";
const PANEL_ANNOUNCE_MODAL_ID = "panel:announce:submit";
const PANEL_ANNOUNCE_TITLE_ID = "panel:announce:title";
const PANEL_ANNOUNCE_MESSAGE_ID = "panel:announce:message";
const PANEL_COLOR = 0x5865f2;

export const data = new SlashCommandBuilder()
  .setName("panel")
  .setDescription("Gere le panneau de controle du staff.")
  .setDMPermission(false)
  .setDefaultMemberPermissions(PermissionFlagsBits.ManageGuild)
  .addSubcommand((subcommand) =>
    subcommand.setName("setup").setDescription("Cree ou repare le panneau dans le salon configure.")
  )
  .addSubcommand((subcommand) =>
    subcommand.setName("refresh").setDescription("Actualise les informations du panneau.")
  )
  .addSubcommand((subcommand) =>
    subcommand.setName("status").setDescription("Verifie la configuration et l'etat du panneau.")
  )
  .addSubcommand((subcommand) =>
    subcommand.setName("remove").setDescription("Supprime le message du panneau enregistre.")
  );

export async function execute(interaction) {
  if (!interaction.inCachedGuild()) {
    await interaction.reply({
      content: "Cette commande doit etre utilisee dans un serveur.",
      flags: MessageFlags.Ephemeral,
    });
    return;
  }

  if (!isPanelStaff(interaction)) {
    await interaction.reply({
      content: "Vous n'avez pas acces au panneau du staff.",
      flags: MessageFlags.Ephemeral,
    });
    return;
  }

  const subcommand = interaction.options.getSubcommand();

  if (subcommand === "status") {
    await sendPanelStatus(interaction);
    return;
  }

  if (subcommand === "remove") {
    await removeStaffPanel(interaction);
    return;
  }

  await interaction.deferReply({ flags: MessageFlags.Ephemeral });

  try {
    const result = await ensureStaffPanel(interaction.guild);
    addAudit({
      guildId: interaction.guildId,
      action: subcommand === "setup" ? "panel_setup" : "panel_refresh",
      actor: interaction.user.tag,
      target: result.channel.id,
    });
    await interaction.editReply(
      result.created
        ? `Le panneau staff a ete cree dans ${result.channel}.`
        : `Le panneau staff a ete actualise dans ${result.channel}.`
    );
  } catch (error) {
    await interaction.editReply(getPanelErrorMessage(error));
  }
}

export async function handlePanelInteraction(interaction) {
  if (
    !interaction.isButton() ||
    ![
      PANEL_REFRESH_BUTTON_ID,
      PANEL_TICKETS_BUTTON_ID,
      PANEL_INCIDENTS_BUTTON_ID,
      PANEL_REPORT_BUTTON_ID,
      PANEL_ANNOUNCE_BUTTON_ID,
    ].includes(interaction.customId)
  ) {
    return false;
  }

  if (!interaction.inCachedGuild() || !isPanelStaff(interaction)) {
    await interaction.reply({
      content: "Vous n'avez pas acces a cette action staff.",
      flags: MessageFlags.Ephemeral,
    });
    return true;
  }

  if (interaction.customId === PANEL_ANNOUNCE_BUTTON_ID) {
    await interaction.showModal(buildAnnouncementModal());
    return true;
  }

  if (interaction.customId === PANEL_REFRESH_BUTTON_ID) {
    await interaction.deferReply({ flags: MessageFlags.Ephemeral });
    try {
      await ensureStaffPanel(interaction.guild);
      addAudit({
        guildId: interaction.guildId,
        action: "panel_refresh",
        actor: interaction.user.tag,
        target: interaction.channelId,
      });
      await interaction.editReply("Le panneau a ete actualise.");
    } catch (error) {
      await interaction.editReply(getPanelErrorMessage(error));
    }
    return true;
  }

  const embed =
    interaction.customId === PANEL_TICKETS_BUTTON_ID
      ? buildTicketsEmbed(interaction.guildId)
      : interaction.customId === PANEL_INCIDENTS_BUTTON_ID
        ? buildIncidentsEmbed(interaction.guildId)
        : buildStaffReportEmbed(interaction.guild);

  const action = interaction.customId.split(":").at(-1);
  addAudit({
    guildId: interaction.guildId,
    action: `panel_${action}`,
    actor: interaction.user.tag,
  });

  await interaction.reply({
    embeds: [embed],
    flags: MessageFlags.Ephemeral,
  });
  return true;
}

export async function handlePanelModalSubmit(interaction) {
  if (!interaction.isModalSubmit() || interaction.customId !== PANEL_ANNOUNCE_MODAL_ID) {
    return false;
  }

  if (!interaction.inCachedGuild() || !isPanelStaff(interaction)) {
    await interaction.reply({
      content: "Vous n'avez pas acces a cette action staff.",
      flags: MessageFlags.Ephemeral,
    });
    return true;
  }

  const channel = await fetchPanelChannel(interaction.guild);
  if (!channel) {
    await interaction.reply({
      content: "Le salon du panneau est introuvable. Verifiez `DISCORD_PANEL_CHANNEL_ID`.",
      flags: MessageFlags.Ephemeral,
    });
    return true;
  }

  const permissionError = getPanelChannelPermissionError(channel, interaction.guild.members.me);
  if (permissionError) {
    await interaction.reply({
      content: permissionError,
      flags: MessageFlags.Ephemeral,
    });
    return true;
  }

  const title = interaction.fields.getTextInputValue(PANEL_ANNOUNCE_TITLE_ID).trim();
  const message = interaction.fields.getTextInputValue(PANEL_ANNOUNCE_MESSAGE_ID).trim();
  const embed = new EmbedBuilder()
    .setColor(PANEL_COLOR)
    .setTitle(title)
    .setDescription(message)
    .setAuthor({
      name: `${interaction.guild.name} • Communication staff`,
      iconURL: interaction.guild.iconURL() ?? undefined,
    })
    .setFooter({ text: `Publie par ${interaction.user.tag}` })
    .setTimestamp();

  await channel.send({
    embeds: [embed],
    allowedMentions: { parse: [] },
  });

  addAudit({
    guildId: interaction.guildId,
    action: "panel_announce",
    actor: interaction.user.tag,
    target: channel.id,
    reason: title,
  });

  await interaction.reply({
    content: `L'annonce staff a ete publiee dans ${channel}.`,
    flags: MessageFlags.Ephemeral,
  });
  return true;
}

export async function ensureStaffPanel(guild) {
  const channel = await fetchPanelChannel(guild);
  if (!channel) {
    throw new Error(
      "Le salon du panneau est introuvable. Configurez `DISCORD_PANEL_CHANNEL_ID` avec un salon textuel."
    );
  }

  const permissionError = getPanelChannelPermissionError(channel, guild.members.me);
  if (permissionError) {
    throw new Error(permissionError);
  }

  const settings = getPanelSettings(guild.id);
  const fetchedMessage =
    settings.messageId && settings.channelId === channel.id
      ? await channel.messages.fetch(settings.messageId).catch(() => null)
      : null;
  const existingMessage =
    fetchedMessage?.author.id === guild.client.user.id ? fetchedMessage : null;
  const payload = buildPanelPayload(guild);

  if (existingMessage) {
    await existingMessage.edit(payload);
    updatePanelSettings(guild.id, {
      channelId: channel.id,
      messageId: existingMessage.id,
      updatedAt: new Date().toISOString(),
    });
    return { channel, message: existingMessage, created: false };
  }

  if (settings.channelId && settings.channelId !== channel.id && settings.messageId) {
    const previousChannel = await guild.channels.fetch(settings.channelId).catch(() => null);
    const previousMessage = previousChannel?.isTextBased()
      ? await previousChannel.messages.fetch(settings.messageId).catch(() => null)
      : null;

    if (previousMessage?.author.id === guild.client.user.id) {
      await previousMessage.delete().catch(() => {});
    }
  }

  const message = await channel.send(payload);
  updatePanelSettings(guild.id, {
    channelId: channel.id,
    messageId: message.id,
    updatedAt: new Date().toISOString(),
  });
  return { channel, message, created: true };
}

async function sendPanelStatus(interaction) {
  const channel = await fetchPanelChannel(interaction.guild);
  const settings = getPanelSettings(interaction.guildId);
  const fetchedMessage =
    channel && settings.messageId && settings.channelId === channel.id
      ? await channel.messages.fetch(settings.messageId).catch(() => null)
      : null;
  const message = fetchedMessage?.author.id === interaction.client.user.id ? fetchedMessage : null;
  const permissionError = channel
    ? getPanelChannelPermissionError(channel, interaction.guild.members.me)
    : null;

  await interaction.reply({
    embeds: [
      new EmbedBuilder()
        .setColor(channel && message && !permissionError ? 0x57f287 : 0xfee75c)
        .setTitle("Etat du panneau staff")
        .addFields(
          {
            name: "Configuration",
            value: env.panelChannelId
              ? `Salon configure: <#${env.panelChannelId}>`
              : "`DISCORD_PANEL_CHANNEL_ID` n'est pas configure.",
          },
          {
            name: "Message persistant",
            value: message
              ? `[Ouvrir le panneau](${message.url})`
              : "Aucun message valide n'est actuellement enregistre.",
          },
          {
            name: "Permissions",
            value:
              permissionError ??
              (channel ? "Permissions requises disponibles." : "Salon invalide."),
          },
          {
            name: "Derniere actualisation",
            value: settings.updatedAt
              ? `<t:${Math.floor(new Date(settings.updatedAt).getTime() / 1000)}:R>`
              : "Jamais",
          }
        )
        .setTimestamp(),
    ],
    flags: MessageFlags.Ephemeral,
  });
}

async function removeStaffPanel(interaction) {
  const settings = getPanelSettings(interaction.guildId);
  const channel = settings.channelId
    ? await interaction.guild.channels.fetch(settings.channelId).catch(() => null)
    : null;
  const fetchedMessage =
    channel?.isTextBased() && settings.messageId
      ? await channel.messages.fetch(settings.messageId).catch(() => null)
      : null;
  const message = fetchedMessage?.author.id === interaction.client.user.id ? fetchedMessage : null;

  if (message) {
    await message.delete();
  }

  resetPanelSettings(interaction.guildId);
  addAudit({
    guildId: interaction.guildId,
    action: "panel_remove",
    actor: interaction.user.tag,
    target: settings.channelId,
  });

  await interaction.reply({
    content: message
      ? "Le panneau staff a ete supprime."
      : "La configuration du panneau a ete reinitialisee.",
    flags: MessageFlags.Ephemeral,
  });
}

function buildPanelPayload(guild) {
  const guildTickets = getGuildTickets(guild.id);
  const openIncidents = getGuildIncidents(guild.id).filter(
    (incident) => incident.status !== "resolved"
  );
  const unclaimedTickets = guildTickets.filter(([, ticket]) => !ticket.claimedBy);
  const lastAuditEntry = auditLog.find((entry) => entry.guildId === guild.id);

  const overviewEmbed = new EmbedBuilder()
    .setColor(PANEL_COLOR)
    .setAuthor({
      name: `${guild.name} • Centre de controle staff`,
      iconURL: guild.iconURL() ?? undefined,
    })
    .setTitle("Pilotage et supervision")
    .setDescription(
      "Ce message centralise les informations operationnelles et les actions rapides reservees au staff."
    )
    .addFields(
      { name: "Membres", value: String(guild.memberCount), inline: true },
      { name: "Tickets ouverts", value: String(guildTickets.length), inline: true },
      { name: "Non assignes", value: String(unclaimedTickets.length), inline: true },
      { name: "Incidents actifs", value: String(openIncidents.length), inline: true },
      {
        name: "Derniere action",
        value: lastAuditEntry
          ? `${lastAuditEntry.action} par ${lastAuditEntry.actor} • <t:${Math.floor(lastAuditEntry.at.getTime() / 1000)}:R>`
          : "Aucune action enregistree.",
        inline: false,
      }
    )
    .setFooter({ text: "Utilisez les boutons ci-dessous pour ouvrir une vue privee." })
    .setTimestamp();

  const operationsEmbed = new EmbedBuilder()
    .setColor(0x2b2d31)
    .setTitle("Etat operationnel")
    .addFields(
      {
        name: "Support",
        value: guildTickets.length
          ? `${guildTickets.length} ticket(s), dont ${unclaimedTickets.length} sans responsable.`
          : "Aucun ticket ouvert.",
        inline: true,
      },
      {
        name: "Incidents",
        value: openIncidents.length
          ? openIncidents
              .slice(0, 3)
              .map((incident) => `• ${incident.id} — ${incident.severity}`)
              .join("\n")
          : "Aucun incident actif.",
        inline: true,
      }
    );

  const actions = new ActionRowBuilder().addComponents(
    new ButtonBuilder()
      .setCustomId(PANEL_REFRESH_BUTTON_ID)
      .setLabel("Actualiser")
      .setStyle(ButtonStyle.Primary),
    new ButtonBuilder()
      .setCustomId(PANEL_TICKETS_BUTTON_ID)
      .setLabel("Tickets")
      .setStyle(ButtonStyle.Secondary),
    new ButtonBuilder()
      .setCustomId(PANEL_INCIDENTS_BUTTON_ID)
      .setLabel("Incidents")
      .setStyle(ButtonStyle.Secondary),
    new ButtonBuilder()
      .setCustomId(PANEL_REPORT_BUTTON_ID)
      .setLabel("Rapport")
      .setStyle(ButtonStyle.Secondary),
    new ButtonBuilder()
      .setCustomId(PANEL_ANNOUNCE_BUTTON_ID)
      .setLabel("Annoncer")
      .setStyle(ButtonStyle.Success)
  );

  return {
    embeds: [overviewEmbed, operationsEmbed],
    components: [actions],
    allowedMentions: { parse: [] },
  };
}

function buildTicketsEmbed(guildId) {
  const guildTickets = getGuildTickets(guildId);
  const lines = guildTickets.slice(0, 15).map(([channelId, ticket]) => {
    const owner = ticket.owner ? `<@${ticket.owner}>` : "inconnu";
    const assignee = ticket.claimedBy ? `<@${ticket.claimedBy}>` : "non assigne";
    return `• <#${channelId}> — ${owner} — ${assignee}${ticket.topic ? ` — ${ticket.topic}` : ""}`;
  });

  return new EmbedBuilder()
    .setColor(PANEL_COLOR)
    .setTitle(`Tickets ouverts • ${guildTickets.length}`)
    .setDescription(lines.join("\n") || "Aucun ticket ouvert.")
    .setFooter({
      text:
        guildTickets.length > 15
          ? `${guildTickets.length - 15} ticket(s) supplementaire(s) non affiche(s).`
          : "Vue privee reservee au staff.",
    })
    .setTimestamp();
}

function buildIncidentsEmbed(guildId) {
  const guildIncidents = getGuildIncidents(guildId);
  const lines = guildIncidents
    .slice(0, 15)
    .map(
      (incident) =>
        `• **${incident.id}** — ${incident.severity} — ${incident.status} — ${incident.title}`
    );

  return new EmbedBuilder()
    .setColor(
      guildIncidents.some((incident) => incident.status !== "resolved") ? 0xed4245 : 0x57f287
    )
    .setTitle(`Incidents • ${guildIncidents.length}`)
    .setDescription(lines.join("\n") || "Aucun incident enregistre.")
    .setFooter({ text: "Utilisez /incident pour creer ou resoudre un incident." })
    .setTimestamp();
}

function buildStaffReportEmbed(guild) {
  const guildTickets = getGuildTickets(guild.id);
  const guildIncidents = getGuildIncidents(guild.id);
  const warningCount = [...warnings.entries()]
    .filter(([key]) => key.startsWith(`${guild.id}:`))
    .reduce((total, [, entries]) => total + entries.length, 0);
  const recentAudit = auditLog
    .filter((entry) => entry.guildId === guild.id)
    .slice(0, 5)
    .map(
      (entry) =>
        `• ${entry.action} — ${entry.actor} — <t:${Math.floor(entry.at.getTime() / 1000)}:R>`
    )
    .join("\n");

  return new EmbedBuilder()
    .setColor(PANEL_COLOR)
    .setTitle(`Rapport staff • ${guild.name}`)
    .addFields(
      {
        name: "Vue d'ensemble",
        value: [
          `Membres: ${guild.memberCount}`,
          `Tickets ouverts: ${guildTickets.length}`,
          `Tickets non assignes: ${guildTickets.filter(([, ticket]) => !ticket.claimedBy).length}`,
          `Incidents actifs: ${guildIncidents.filter((incident) => incident.status !== "resolved").length}`,
          `Warnings stockes: ${warningCount}`,
        ].join("\n"),
      },
      {
        name: "Dernieres actions",
        value: recentAudit || "Aucune action recente.",
      }
    )
    .setTimestamp();
}

function buildAnnouncementModal() {
  const titleInput = new TextInputBuilder()
    .setCustomId(PANEL_ANNOUNCE_TITLE_ID)
    .setLabel("Titre")
    .setStyle(TextInputStyle.Short)
    .setMinLength(3)
    .setMaxLength(100)
    .setRequired(true);
  const messageInput = new TextInputBuilder()
    .setCustomId(PANEL_ANNOUNCE_MESSAGE_ID)
    .setLabel("Message")
    .setStyle(TextInputStyle.Paragraph)
    .setMinLength(3)
    .setMaxLength(2000)
    .setRequired(true);

  return new ModalBuilder()
    .setCustomId(PANEL_ANNOUNCE_MODAL_ID)
    .setTitle("Publier une annonce staff")
    .addComponents(
      new ActionRowBuilder().addComponents(titleInput),
      new ActionRowBuilder().addComponents(messageInput)
    );
}

function getGuildTickets(guildId) {
  return [...tickets.entries()].filter(([, ticket]) => ticket.guildId === guildId);
}

function getGuildIncidents(guildId) {
  return [...incidents.values()].filter(
    (incident) => !incident.guildId || incident.guildId === guildId
  );
}

async function fetchPanelChannel(guild) {
  if (!env.panelChannelId) {
    return null;
  }

  const channel =
    guild.channels.cache.get(env.panelChannelId) ??
    (await guild.channels.fetch(env.panelChannelId).catch(() => null));

  if (
    !channel ||
    !channel.isTextBased() ||
    channel.isDMBased() ||
    ![ChannelType.GuildText, ChannelType.GuildAnnouncement].includes(channel.type)
  ) {
    return null;
  }

  return channel;
}

function getPanelChannelPermissionError(channel, botMember) {
  const permissions = botMember ? channel.permissionsFor(botMember) : null;
  const requiredPermissions = [
    [PermissionFlagsBits.ViewChannel, "Voir le salon"],
    [PermissionFlagsBits.SendMessages, "Envoyer des messages"],
    [PermissionFlagsBits.ReadMessageHistory, "Voir les anciens messages"],
    [PermissionFlagsBits.EmbedLinks, "Integrer des liens"],
  ];
  const missingPermissions = requiredPermissions
    .filter(([permission]) => !permissions?.has(permission))
    .map(([, label]) => label);

  if (missingPermissions.length === 0) {
    return null;
  }

  return `Permissions manquantes dans ${channel}: ${missingPermissions.join(", ")}.`;
}

function isPanelStaff(interaction) {
  const hasStaffRole = Boolean(
    env.modoRoleId && interaction.member?.roles?.cache?.has(env.modoRoleId)
  );
  const hasStaffPermission =
    interaction.memberPermissions?.has(PermissionFlagsBits.ManageGuild) ||
    interaction.memberPermissions?.has(PermissionFlagsBits.ManageChannels) ||
    interaction.memberPermissions?.has(PermissionFlagsBits.ModerateMembers);

  return hasStaffRole || Boolean(hasStaffPermission);
}

function getPanelErrorMessage(error) {
  return error instanceof Error
    ? error.message
    : "Impossible de creer ou d'actualiser le panneau staff.";
}
