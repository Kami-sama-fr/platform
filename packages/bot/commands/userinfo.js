import { EmbedBuilder, SlashCommandBuilder } from "discord.js";

const STATUS_EMOJIS = {
  online: "🟢",
  idle: "🟡",
  dnd: "🔴",
  offline: "⚫",
};

export const data = new SlashCommandBuilder()
  .setName("userinfo")
  .setDescription("Affiche les informations detaillees d'un utilisateur.")
  .addUserOption((o) =>
    o.setName("user").setDescription("Utilisateur cible (defaut: vous)").setRequired(false)
  );

export async function execute(interaction) {
  if (!interaction.inCachedGuild()) {
    await interaction.reply({ content: "Cette commande doit etre utilisee dans un serveur.", ephemeral: true });
    return;
  }

  const user = interaction.options.getUser("user") ?? interaction.user;
  const member = await interaction.guild.members.fetch(user.id).catch(() => null);

  const created = Math.floor(user.createdTimestamp / 1000);
  const joined = member?.joinedAt ? Math.floor(member.joinedAt.getTime() / 1000) : null;
  const status = member?.presence?.status ?? "offline";
  const isBot = user.bot ? "Oui" : "Non";
  const isBoosting = member?.premiumSince ? `<t:${Math.floor(member.premiumSince.getTime() / 1000)}:F>` : "Non";
  const roles = member?.roles.cache.filter((r) => r.id !== interaction.guild.id).sort((a, b) => b.position - a.position);
  const rolesList = roles?.size > 0 ? roles.map((r) => r.toString()).join(" ") : "Aucun";
  const topRoleColor = member?.roles.highest.color ?? 0x5865f2;

  const embed = new EmbedBuilder()
    .setColor(topRoleColor)
    .setAuthor({ name: user.tag, iconURL: user.displayAvatarURL() })
    .setThumbnail(user.displayAvatarURL({ size: 1024 }))
    .addFields(
      { name: "Informations generales", value: [
        `**ID:** \`${user.id}\``,
        `**Affichage:** ${member?.displayName ?? user.username}`,
        `**Bot:** ${isBot}`,
        `**Statut:** ${STATUS_EMOJIS[status]} ${status}`,
      ].join("\n"), inline: false },
      { name: "Dates", value: [
        `**Compte cree:** <t:${created}:F> (<t:${created}:R>)`,
        joined ? `**A rejoint:** <t:${joined}:F> (<t:${joined}:R>)` : "**A rejoint:** Inconnu",
      ].join("\n"), inline: false },
      { name: "Boosting", value: isBoosting, inline: true },
      { name: "Roles", value: roles?.size > 0 ? `${roles.size} role(s)\n${rolesList}` : "Aucun role", inline: false },
    )
    .setFooter({ text: `ID: ${user.id}` })
    .setTimestamp();

  if (user.banner) {
    embed.setImage(user.bannerURL({ size: 1024 }));
  }

  await interaction.reply({ embeds: [embed], ephemeral: true });
}
