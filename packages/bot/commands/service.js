import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  ChannelType,
  EmbedBuilder,
  PermissionFlagsBits,
  SlashCommandBuilder,
} from "discord.js";
import { addAudit, getNotificationChannels, getNotificationTypes, removeNotificationChannel, setNotificationChannel } from "../utils/store.js";

const SUMMARY_BUTTON_ID = "service:refresh";

const LABELS = {
  server_info: { emoji: "🖥️", label: "Server Info", description: "Notifications sur l'etat du serveur et des services principaux." },
  infra_info: { emoji: "⚙️", label: "Infra Info", description: "Alertes et informations sur l'infrastructure technique." },
  social_network: { emoji: "🌐", label: "Social Network", description: "Activite et notifications des reseaux sociaux." },
  status: { emoji: "📊", label: "Status", description: "Rapports de statut et bilans periodiques." },
};

const choices = getNotificationTypes().map((t) => ({
  name: LABELS[t]?.label ?? t,
  value: t,
}));

export const data = new SlashCommandBuilder()
  .setName("service")
  .setDescription("Configure les modules de notifications pour les services.")
  .setDMPermission(false)
  .setDefaultMemberPermissions(PermissionFlagsBits.ManageGuild)
  .addSubcommand((sub) =>
    sub
      .setName("setup")
      .setDescription("Associe un type de notification a un salon.")
      .addStringOption((o) =>
        o.setName("type").setDescription("Type de notification").setRequired(true).addChoices(...choices)
      )
      .addChannelOption((o) =>
        o.setName("channel").setDescription("Salon de notification").setRequired(true).addChannelTypes(ChannelType.GuildText, ChannelType.GuildAnnouncement)
      )
  )
  .addSubcommand((sub) =>
    sub
      .setName("remove")
      .setDescription("Supprime la configuration d'un type de notification.")
      .addStringOption((o) =>
        o.setName("type").setDescription("Type de notification").setRequired(true).addChoices(...choices)
      )
  )
  .addSubcommand((sub) =>
    sub.setName("status").setDescription("Affiche la configuration actuelle des notifications.")
  )
  .addSubcommand((sub) =>
    sub.setName("summary").setDescription("Publie un tableau de bord resume des modules de notification.")
  );

export async function execute(interaction) {
  if (!interaction.inCachedGuild()) {
    await interaction.reply({ content: "Cette commande doit etre utilisee dans un serveur.", ephemeral: true });
    return;
  }

  const sub = interaction.options.getSubcommand();

  if (sub === "setup") {
    await handleSetup(interaction);
  } else if (sub === "remove") {
    await handleRemove(interaction);
  } else if (sub === "status") {
    await handleStatus(interaction);
  } else if (sub === "summary") {
    await handleSummary(interaction);
  }
}

async function handleSetup(interaction) {
  const type = interaction.options.getString("type", true);
  const channel = interaction.options.getChannel("channel", true);

  const botMember = interaction.guild.members.me;
  const perms = channel.permissionsFor(botMember);
  if (!perms?.has(PermissionFlagsBits.ViewChannel) || !perms.has(PermissionFlagsBits.SendMessages)) {
    await interaction.reply({ content: `Le bot n'a pas les permissions necessaires dans ${channel}.`, ephemeral: true });
    return;
  }

  setNotificationChannel(interaction.guildId, type, channel.id);

  addAudit({
    guildId: interaction.guildId,
    action: "service_setup",
    actor: interaction.user.tag,
    target: channel.id,
    reason: type,
  });

  const embed = new EmbedBuilder()
    .setColor(0x57f287)
    .setTitle(`${LABELS[type].emoji} Module ${LABELS[type].label} active`)
    .setDescription(
      [
        `Ce salon est desormais configure pour recevoir les notifications **${LABELS[type].label}**.`,
        "",
        `${LABELS[type].description}`,
      ].join("\n")
    )
    .setFooter({ text: `Configure par ${interaction.user.tag}` })
    .setTimestamp();

  await channel.send({ embeds: [embed] });

  await interaction.reply({
    content: `Notifications **${LABELS[type].label}** configurees dans ${channel}. Un message de confirmation a ete publie.`,
    ephemeral: true,
  });
}

async function handleRemove(interaction) {
  const type = interaction.options.getString("type", true);
  const channels = getNotificationChannels(interaction.guildId);

  if (!channels[type]) {
    await interaction.reply({ content: `Aucune configuration trouvee pour **${LABELS[type]?.label ?? type}**.`, ephemeral: true });
    return;
  }

  removeNotificationChannel(interaction.guildId, type);

  addAudit({
    guildId: interaction.guildId,
    action: "service_remove",
    actor: interaction.user.tag,
    reason: type,
  });

  await interaction.reply({
    content: `Notifications **${LABELS[type]?.label ?? type}** supprimees.`,
    ephemeral: true,
  });
}

async function handleStatus(interaction) {
  const channels = getNotificationChannels(interaction.guildId);
  const fields = getNotificationTypes().map((type) => {
    const info = LABELS[type];
    const channelId = channels[type];
    return {
      name: `${info.emoji} ${info.label}`,
      value: channelId ? `<#${channelId}>` : "❌ Non configure",
      inline: true,
    };
  });

  const embed = new EmbedBuilder()
    .setColor(0x5865f2)
    .setTitle("Configuration des notifications")
    .setDescription("Salons associes a chaque module de notification.")
    .addFields(...fields)
    .setFooter({ text: interaction.guild.name })
    .setTimestamp();

  await interaction.reply({ embeds: [embed], ephemeral: true });
}

function buildSummaryEmbed(guild, channels) {
  const total = getNotificationTypes().length;
  const configured = getNotificationTypes().filter((t) => channels[t]).length;

  const fields = getNotificationTypes().map((type) => {
    const info = LABELS[type];
    const channelId = channels[type];
    return {
      name: `${channelId ? "✅" : "⛔"} ${info.emoji} ${info.label}`,
      value: channelId
        ? `Salon: <#${channelId}>\n${info.description}`
        : `${info.description}\n\n_Aucun salon configure._`,
      inline: false,
    };
  });

  return new EmbedBuilder()
    .setColor(configured === total ? 0x57f287 : 0xfee75c)
    .setTitle("Tableau de bord des services")
    .setDescription(
      [
        `Modules de notification pour **${guild.name}**.`,
        `${configured}/${total} modules actifs.`,
      ].join("\n")
    )
    .addFields(...fields)
    .setFooter({ text: "Cliquez sur Actualiser pour mettre a jour les donnees." })
    .setTimestamp();
}

function buildSummaryComponents() {
  return [
    new ActionRowBuilder().addComponents(
      new ButtonBuilder()
        .setCustomId(SUMMARY_BUTTON_ID)
        .setLabel("Actualiser")
        .setStyle(ButtonStyle.Primary)
        .setEmoji("🔄")
    ),
  ];
}

async function handleSummary(interaction) {
  const channels = getNotificationChannels(interaction.guildId);

  if (Object.keys(channels).length === 0) {
    await interaction.reply({
      content: "Aucun module de notification configure. Utilisez `/service setup` pour configurer un salon.",
      ephemeral: true,
    });
    return;
  }

  await interaction.reply({
    embeds: [buildSummaryEmbed(interaction.guild, channels)],
    components: buildSummaryComponents(),
  });
}

export async function handleServiceButtonInteraction(interaction) {
  if (!interaction.isButton() || interaction.customId !== SUMMARY_BUTTON_ID) {
    return false;
  }

  if (!interaction.inCachedGuild()) {
    await interaction.reply({ content: "Cette action doit etre utilisee depuis un serveur.", ephemeral: true });
    return true;
  }

  const channels = getNotificationChannels(interaction.guildId);

  await interaction.update({
    embeds: [buildSummaryEmbed(interaction.guild, channels)],
    components: buildSummaryComponents(),
  });

  return true;
}
