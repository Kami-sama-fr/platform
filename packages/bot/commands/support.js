import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  EmbedBuilder,
  MessageFlags,
  SlashCommandBuilder,
  StringSelectMenuBuilder,
} from "discord.js";
import {
  claimSupportTicket,
  closeSupportTicket,
  createSupportTicket,
  getSupportTicket,
  isSupportStaff,
} from "../utils/support-tickets.js";
import { env } from "../config/env.js";

const OPEN_TICKET_BUTTON_ID = "support:open-ticket";
const CLAIM_TICKET_BUTTON_ID = "support:claim-ticket";
const CLOSE_TICKET_BUTTON_ID = "support:close-ticket";
const LANGUAGE_SELECT_ID = "support:language";
const TOPIC_SELECT_ID = "support:topic";
const CONFIRM_TICKET_BUTTON_ID = "support:confirm-ticket";

const supportRequestDrafts = new Map();

const SUPPORT_LANGUAGES = [
  { value: "fr", label: "Francais", emoji: "🇫🇷", description: "Support en francais" },
  { value: "en", label: "English", emoji: "🇬🇧", description: "Support in English" },
];

const SUPPORT_TOPICS = [
  { value: "account", label: "Compte et acces", emoji: "🔐", description: "Connexion, verification, permissions" },
  { value: "billing", label: "Facturation", emoji: "💳", description: "Paiement, abonnement, facture" },
  { value: "technical", label: "Probleme technique", emoji: "🛠️", description: "Bug, erreur, comportement anormal" },
  { value: "community", label: "Moderation et communaute", emoji: "🧭", description: "Signalement, regles, organisation" },
  { value: "other", label: "Autre demande", emoji: "📩", description: "Besoin specifique non liste" },
];

export const data = new SlashCommandBuilder()
  .setName("support")
  .setDescription("Publie le panneau de support ou ferme un ticket.")
  .addSubcommand((subcommand) =>
    subcommand
      .setName("setup")
      .setDescription("Publie le message de support avec bouton d'ouverture."),
  )
  .addSubcommand((subcommand) =>
    subcommand
      .setName("close")
      .setDescription("Ferme le ticket du salon actuel."),
  );

export async function execute(interaction) {
  const subcommand = interaction.options.getSubcommand();

  if (subcommand === "setup") {
    if (!isSupportStaff(interaction)) {
      await interaction.reply({
        content: "Vous devez avoir la permission de gérer les salons pour publier le panneau de support.",
        ephemeral: true,
      });
      return;
    }

    await sendSupportPanel(interaction.channel);

    await interaction.reply({
      content: "Le panneau de support a été publié dans ce salon.",
      ephemeral: true,
    });
    return;
  }

  const ticket = getSupportTicket(interaction.channelId);

  if (!ticket) {
    await interaction.reply({
      content: "Ce salon n’est pas un ticket géré par le bot.",
      ephemeral: true,
    });
    return;
  }

  if (ticket.owner !== interaction.user.id && !isSupportStaff(interaction)) {
    await interaction.reply({
      content: "Seul l'auteur du ticket ou un modérateur peut le fermer.",
      ephemeral: true,
    });
    return;
  }

  await interaction.reply({
    content: "Fermeture du ticket…",
    ephemeral: true,
  });

  closeSupportTicket(interaction.channelId);
  setTimeout(() => interaction.channel.delete().catch(() => {}), 1000);
}

export async function sendSupportPanel(channel) {
  const guildName = channel.guild?.name ?? "Sky Genesis Enterprise";
  const guildIconUrl = channel.guild?.iconURL?.() ?? null;

  const overviewEmbed = new EmbedBuilder()
    .setTitle(`Centre de support • ${guildName}`)
    .setDescription(
      [
        "Bienvenue dans l'espace d'assistance officiel.",
        "",
        "Utilisez ce panneau pour preparer une demande de support claire et l'envoyer a l'equipe dans un ticket prive.",
      ].join("\n"),
    )
    .addFields(
      {
        name: "Demandes prises en charge",
        value: [
          "• Problemes d'acces ou de compte",
          "• Support technique sur les services de l'ecosysteme",
          "• Signalement de bug ou comportement anormal",
          "• Question sur les procedures ou permissions",
        ].join("\n"),
        inline: false,
      },
      {
        name: "Parcours de demande",
        value: [
          "1. Cliquez sur **Demarrer une demande**.",
          "2. Choisissez votre langue et le sujet correspondant.",
          "3. Validez pour creer votre ticket prive.",
          "4. Decrivez votre probleme avec un maximum de contexte.",
          "5. Un membre du staff prendra votre demande en charge.",
        ].join("\n"),
        inline: false,
      },
    )
    .setColor(0x5865f2)
    .setThumbnail(guildIconUrl)
    .setFooter({
      text: "Support Sky Genesis Enterprise",
    })
    .setTimestamp();

  const guidanceEmbed = new EmbedBuilder()
    .setTitle("Avant de creer votre ticket")
    .setDescription(
      [
        "Plus votre demande est precise des le depart, plus le traitement sera rapide et pertinent.",
      ].join("\n"),
    )
    .addFields(
      {
        name: "A preparer",
        value: [
          "• Le service concerne",
          "• Les etapes deja effectuees",
          "• Le message d'erreur exact",
          "• Des captures d'ecran si necessaire",
        ].join("\n"),
        inline: true,
      },
      {
        name: "Recommandations",
        value: [
          "• Un seul ticket par sujet",
          "• Pas de spam ou relances multiples",
          "• Ne partagez jamais de donnees sensibles",
          "• Fermez le ticket une fois resolu",
        ].join("\n"),
        inline: true,
      },
    )
    .setColor(0x57f287);

  const row = new ActionRowBuilder().addComponents(
    new ButtonBuilder()
      .setCustomId(OPEN_TICKET_BUTTON_ID)
      .setLabel("Demarrer une demande")
      .setStyle(ButtonStyle.Primary),
  );

  await channel.send({
    embeds: [overviewEmbed, guidanceEmbed],
    components: [row],
  });
}

export async function handleSupportButtonInteraction(interaction) {
  if (interaction.isButton() && interaction.customId === CLAIM_TICKET_BUTTON_ID) {
    const ticket = getSupportTicket(interaction.channelId);
    if (!ticket) {
      await interaction.reply({
        content: "Ce salon n’est pas un ticket géré par le bot.",
        flags: MessageFlags.Ephemeral,
      });
      return true;
    }

    if (!isSupportStaff(interaction)) {
      await interaction.reply({
        content: "Seul un membre du staff peut prendre ce ticket.",
        flags: MessageFlags.Ephemeral,
      });
      return true;
    }

    claimSupportTicket(interaction.channelId, interaction.user.id, interaction.user.tag);
    await interaction.reply({
      content: `${interaction.user} a pris en charge ce ticket.`,
    });
    return true;
  }

  if (interaction.isButton() && interaction.customId === CLOSE_TICKET_BUTTON_ID) {
    const ticket = getSupportTicket(interaction.channelId);
    if (!ticket) {
      await interaction.reply({
        content: "Ce salon n’est pas un ticket géré par le bot.",
        flags: MessageFlags.Ephemeral,
      });
      return true;
    }

    if (ticket.owner !== interaction.user.id && !isSupportStaff(interaction)) {
      await interaction.reply({
        content: "Seul l'auteur du ticket ou un modérateur peut le fermer.",
        flags: MessageFlags.Ephemeral,
      });
      return true;
    }

    await interaction.reply({
      content: "Fermeture du ticket…",
      flags: MessageFlags.Ephemeral,
    });

    closeSupportTicket(interaction.channelId);
    setTimeout(() => interaction.channel.delete().catch(() => {}), 1000);
    return true;
  }

  if (interaction.isButton() && interaction.customId === OPEN_TICKET_BUTTON_ID) {
    resetSupportDraft(interaction.user.id);
    await interaction.reply({
      embeds: [buildSupportRequestEmbed(interaction.user.id)],
      components: buildSupportRequestComponents(interaction.user.id),
      flags: MessageFlags.Ephemeral,
    });
    return true;
  }

  if (interaction.isStringSelectMenu() && interaction.customId === LANGUAGE_SELECT_ID) {
    saveSupportDraft(interaction.user.id, { language: interaction.values[0] });
    await interaction.update({
      embeds: [buildSupportRequestEmbed(interaction.user.id)],
      components: buildSupportRequestComponents(interaction.user.id),
    });
    return true;
  }

  if (interaction.isStringSelectMenu() && interaction.customId === TOPIC_SELECT_ID) {
    saveSupportDraft(interaction.user.id, { topic: interaction.values[0] });
    await interaction.update({
      embeds: [buildSupportRequestEmbed(interaction.user.id)],
      components: buildSupportRequestComponents(interaction.user.id),
    });
    return true;
  }

  if (interaction.isButton() && interaction.customId === CONFIRM_TICKET_BUTTON_ID) {
    const draft = supportRequestDrafts.get(interaction.user.id);

    if (!draft?.language || !draft?.topic) {
      await interaction.reply({
        content: "Selectionnez d'abord une langue et un sujet avant de creer le ticket.",
        flags: MessageFlags.Ephemeral,
      });
      return true;
    }

    const channel = await createSupportTicket(interaction, draft);
    resetSupportDraft(interaction.user.id);

    await channel.send({
      content: buildSupportTeamMention(),
      embeds: [
        buildCreatedTicketEmbed(interaction, draft),
      ],
      components: [
        new ActionRowBuilder().addComponents(
          new ButtonBuilder()
            .setCustomId(CLAIM_TICKET_BUTTON_ID)
            .setLabel("Prendre le ticket")
            .setStyle(ButtonStyle.Primary),
          new ButtonBuilder()
          .setCustomId(CLOSE_TICKET_BUTTON_ID)
          .setLabel("Fermer ce ticket")
          .setStyle(ButtonStyle.Danger),
        ),
      ],
      allowedMentions: env.modoRoleId
        ? {
            roles: [env.modoRoleId],
          }
        : undefined,
    });

    await interaction.update({
      content: `Votre ticket a ete cree: ${channel}`,
      embeds: [],
      components: [],
    });
    return true;
  }

  return false;
}

export async function handleSupportModalSubmit() {
  return false;
}

function saveSupportDraft(userId, values) {
  const draft = supportRequestDrafts.get(userId) ?? {};
  supportRequestDrafts.set(userId, {
    ...draft,
    ...values,
    updatedAt: new Date(),
  });
}

function resetSupportDraft(userId) {
  supportRequestDrafts.delete(userId);
}

function buildSupportRequestEmbed(userId) {
  const draft = supportRequestDrafts.get(userId) ?? {};
  const language = getLanguageLabel(draft.language);
  const topic = getTopicLabel(draft.topic);

  return new EmbedBuilder()
    .setTitle("Preparation de votre ticket")
    .setDescription(
      [
        "Choisissez la langue du support et le sujet principal de votre demande.",
        "",
        "Le ticket sera ensuite cree dans un salon prive avec ces informations pour orienter plus vite le staff.",
      ].join("\n"),
    )
    .addFields(
      {
        name: "Langue",
        value: language ?? "Non selectionnee",
        inline: true,
      },
      {
        name: "Sujet",
        value: topic ?? "Non selectionne",
        inline: true,
      },
    )
    .setColor(0x5865f2);
}

function buildSupportRequestComponents(userId) {
  const draft = supportRequestDrafts.get(userId) ?? {};

  const languageRow = new ActionRowBuilder().addComponents(
    new StringSelectMenuBuilder()
      .setCustomId(LANGUAGE_SELECT_ID)
      .setPlaceholder("Choisir une langue")
      .addOptions(
        SUPPORT_LANGUAGES.map((language) => ({
          label: language.label,
          value: language.value,
          emoji: language.emoji,
          description: language.description,
          default: draft.language === language.value,
        })),
      ),
  );

  const topicRow = new ActionRowBuilder().addComponents(
    new StringSelectMenuBuilder()
      .setCustomId(TOPIC_SELECT_ID)
      .setPlaceholder("Choisir un sujet de support")
      .addOptions(
        SUPPORT_TOPICS.map((topic) => ({
          label: topic.label,
          value: topic.value,
          emoji: topic.emoji,
          description: topic.description,
          default: draft.topic === topic.value,
        })),
      ),
  );

  const actionRow = new ActionRowBuilder().addComponents(
    new ButtonBuilder()
      .setCustomId(CONFIRM_TICKET_BUTTON_ID)
      .setLabel("Creer le ticket")
      .setStyle(ButtonStyle.Success)
      .setDisabled(!draft.language || !draft.topic),
  );

  return [languageRow, topicRow, actionRow];
}

function buildCreatedTicketEmbed(interaction, draft) {
  return new EmbedBuilder()
    .setTitle("Nouveau ticket")
    .setDescription("Un membre du staff prendra votre demande en charge des que possible.")
    .addFields(
      { name: "Demandeur", value: `${interaction.user}`, inline: true },
      { name: "Statut", value: "Ouvert", inline: true },
      { name: "Langue", value: getLanguageLabel(draft.language), inline: true },
      { name: "Sujet", value: getTopicLabel(draft.topic), inline: true },
    )
    .setColor(0x57f287);
}

function getLanguageLabel(value) {
  return SUPPORT_LANGUAGES.find((language) => language.value === value)?.label ?? null;
}

function getTopicLabel(value) {
  return SUPPORT_TOPICS.find((topic) => topic.value === value)?.label ?? null;
}

function buildSupportTeamMention() {
  if (!env.modoRoleId) {
    return "L'equipe support a ete notifiee.";
  }

  return `<@&${env.modoRoleId}> nouveau ticket de support cree.`;
}
