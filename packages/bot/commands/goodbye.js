import { ChannelType, PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import {
  buildMemberMessagePayload,
  parseHexColor,
  renderMemberMessage,
  resolveGoodbyeSettings,
} from "../utils/member-messages.js";
import {
  addAudit,
  resetGoodbyeSettings,
  setMemberEventsEnabled,
  updateGoodbyeSettings,
} from "../utils/store.js";

export const data = new SlashCommandBuilder()
  .setName("goodbye")
  .setDescription("Configure les messages de depart des membres.")
  .setDefaultMemberPermissions(PermissionFlagsBits.ManageGuild)
  .addSubcommand((subcommand) =>
    subcommand.setName("enable").setDescription("Active les messages de depart.")
  )
  .addSubcommand((subcommand) =>
    subcommand.setName("disable").setDescription("Desactive les messages de depart.")
  )
  .addSubcommand((subcommand) =>
    subcommand.setName("status").setDescription("Affiche la configuration actuelle.")
  )
  .addSubcommand((subcommand) =>
    subcommand
      .setName("channel")
      .setDescription("Definit le salon des messages de depart.")
      .addChannelOption((option) =>
        option
          .setName("channel")
          .setDescription("Salon cible")
          .addChannelTypes(ChannelType.GuildText, ChannelType.GuildAnnouncement)
          .setRequired(true)
      )
  )
  .addSubcommand((subcommand) =>
    subcommand
      .setName("message")
      .setDescription("Definit le message de depart.")
      .addStringOption((option) =>
        option
          .setName("content")
          .setDescription("Variables: {member}, {username}, {userTag}, {server}, {memberCount}")
          .setMaxLength(2000)
          .setRequired(true)
      )
  )
  .addSubcommand((subcommand) =>
    subcommand.setName("preview").setDescription("Affiche un apercu du message.")
  )
  .addSubcommand((subcommand) =>
    subcommand
      .setName("bots")
      .setDescription("Configure les messages de depart des bots.")
      .addBooleanOption((option) =>
        option
          .setName("ignore")
          .setDescription("Ignorer les bots lors de leur depart")
          .setRequired(true)
      )
  )
  .addSubcommand((subcommand) =>
    subcommand
      .setName("format")
      .setDescription("Choisit un message texte ou un embed.")
      .addStringOption((option) =>
        option
          .setName("type")
          .setDescription("Format du message")
          .addChoices({ name: "Texte", value: "text" }, { name: "Embed", value: "embed" })
          .setRequired(true)
      )
      .addStringOption((option) =>
        option
          .setName("color")
          .setDescription("Couleur hexadecimale de l'embed, par exemple #5865F2")
          .setRequired(false)
      )
  )
  .addSubcommand((subcommand) =>
    subcommand.setName("reset").setDescription("Retablit la configuration par defaut.")
  );

export async function execute(interaction) {
  if (!interaction.guild) {
    await interaction.reply({
      content: "Cette commande doit etre utilisee dans un serveur.",
      ephemeral: true,
    });
    return;
  }

  const guildId = interaction.guild.id;
  const subcommand = interaction.options.getSubcommand();

  if (subcommand === "status") {
    const settings = resolveGoodbyeSettings(guildId);
    await interaction.reply({
      content: [
        `**Messages de depart:** ${settings.enabled ? "actives" : "desactives"}`,
        `**Salon:** ${settings.channelId ? `<#${settings.channelId}>` : "Non configure"}`,
        `**Format:** ${settings.format}`,
        `**Bots ignores:** ${settings.ignoreBots ? "oui" : "non"}`,
        `**Message:** ${settings.message}`,
      ].join("\n"),
      ephemeral: true,
    });
    return;
  }

  if (subcommand === "preview") {
    const settings = resolveGoodbyeSettings(guildId);
    const member = await interaction.guild.members.fetch(interaction.user.id);
    const content = renderMemberMessage(settings.message, member);
    await interaction.reply({
      ...buildMemberMessagePayload(content, settings, `Depart de ${interaction.guild.name}`),
      ephemeral: true,
    });
    return;
  }

  if (subcommand === "enable" || subcommand === "disable") {
    const enabled = subcommand === "enable";
    setMemberEventsEnabled(guildId, true);
    updateGoodbyeSettings(guildId, { enabled });
    await replyWithAudit(
      interaction,
      enabled ? "goodbye_enable" : "goodbye_disable",
      `Les messages de depart sont **${enabled ? "actives" : "desactives"}**.`
    );
    return;
  }

  if (subcommand === "channel") {
    const channel = interaction.options.getChannel("channel", true);
    updateGoodbyeSettings(guildId, { channelId: channel.id });
    await replyWithAudit(
      interaction,
      "goodbye_channel",
      `Le salon de depart est maintenant ${channel}.`
    );
    return;
  }

  if (subcommand === "message") {
    updateGoodbyeSettings(guildId, {
      message: interaction.options.getString("content", true),
    });
    await replyWithAudit(interaction, "goodbye_message", "Le message de depart a ete mis a jour.");
    return;
  }

  if (subcommand === "bots") {
    const ignoreBots = interaction.options.getBoolean("ignore", true);
    updateGoodbyeSettings(guildId, { ignoreBots });
    await replyWithAudit(
      interaction,
      "goodbye_bots",
      `Les comptes bots sont maintenant **${ignoreBots ? "ignores" : "annonces"}**.`
    );
    return;
  }

  if (subcommand === "format") {
    const format = interaction.options.getString("type", true);
    const colorInput = interaction.options.getString("color");
    const color = parseHexColor(colorInput);

    if (colorInput && color === null) {
      await interaction.reply({
        content: "La couleur doit utiliser six caracteres hexadecimaux, par exemple `#5865F2`.",
        ephemeral: true,
      });
      return;
    }

    updateGoodbyeSettings(guildId, {
      format,
      ...(color !== null ? { color } : {}),
    });
    await replyWithAudit(interaction, "goodbye_format", `Le format est maintenant **${format}**.`);
    return;
  }

  if (subcommand === "reset") {
    resetGoodbyeSettings(guildId);
    setMemberEventsEnabled(guildId, true);
    await replyWithAudit(
      interaction,
      "goodbye_reset",
      "La configuration de depart utilise de nouveau les valeurs par defaut."
    );
  }
}

async function replyWithAudit(interaction, action, content) {
  addAudit({
    guildId: interaction.guild.id,
    action,
    actor: interaction.user.tag,
  });
  await interaction.reply({ content, ephemeral: true });
}
