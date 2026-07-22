import {
  ChannelType,
  EmbedBuilder,
  PermissionFlagsBits,
  SlashCommandBuilder,
} from "discord.js";

const ANNOUNCE_COLOR = 0x5865f2;
const ROLE_OPTION_NAMES = ["role", "role_2", "role_3"];

export const data = new SlashCommandBuilder()
  .setName("announce")
  .setDescription("Publie une annonce structuree avec mentions de roles.")
  .setDefaultMemberPermissions(PermissionFlagsBits.ManageMessages)
  .addChannelOption((option) =>
    option
      .setName("channel")
      .setDescription("Salon cible")
      .addChannelTypes(ChannelType.GuildText, ChannelType.GuildAnnouncement)
      .setRequired(true),
  )
  .addStringOption((option) =>
    option.setName("title").setDescription("Titre de l'annonce").setRequired(true),
  )
  .addStringOption((option) =>
    option
      .setName("message")
      .setDescription("Contenu detaille de l'annonce")
      .setRequired(true),
  )
  .addStringOption((option) =>
    option
      .setName("summary")
      .setDescription("Resume court affiche en haut de l'embed"),
  )
  .addRoleOption((option) =>
    option.setName("role").setDescription("Premier role a mentionner"),
  )
  .addRoleOption((option) =>
    option.setName("role_2").setDescription("Deuxieme role a mentionner"),
  )
  .addRoleOption((option) =>
    option.setName("role_3").setDescription("Troisieme role a mentionner"),
  );

export async function execute(interaction) {
  const channel = interaction.options.getChannel("channel", true);

  if (!channel.isTextBased() || !("send" in channel)) {
    await interaction.reply({
      content: "Salon incompatible.",
      ephemeral: true,
    });
    return;
  }

  const roles = ROLE_OPTION_NAMES
    .map((name) => interaction.options.getRole(name))
    .filter((role) => role && role.id !== interaction.guildId)
    .filter((role, index, list) => list.findIndex((entry) => entry.id === role.id) === index);

  const title = interaction.options.getString("title", true);
  const summary = interaction.options.getString("summary");
  const message = interaction.options.getString("message", true);
  const guildName = interaction.guild?.name ?? "Serveur";
  const guildIconUrl = interaction.guild?.iconURL?.() ?? null;

  const embed = new EmbedBuilder()
    .setColor(ANNOUNCE_COLOR)
    .setTitle(title)
    .setDescription(summary ? `**${summary}**\n\n${message}` : message)
    .setAuthor({
      name: `${guildName} • Annonce`,
      iconURL: guildIconUrl ?? undefined,
    })
    .setFooter({
      text: `Publie par ${interaction.user.tag}`,
    })
    .setTimestamp();

  if (roles.length > 0) {
    embed.addFields({
      name: "Roles concernes",
      value: roles.map((role) => `${role}`).join(" "),
      inline: false,
    });
  }

  await channel.send({
    content: roles.length > 0 ? roles.map((role) => `${role}`).join(" ") : undefined,
    embeds: [embed],
    allowedMentions: roles.length > 0
      ? {
          roles: roles.map((role) => role.id),
        }
      : undefined,
  });

  await interaction.reply({
    content: `Annonce publiee dans ${channel}.`,
    ephemeral: true,
  });
}
