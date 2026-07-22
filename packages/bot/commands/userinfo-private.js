import { EmbedBuilder, PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import { warnings, xp } from "../utils/store.js";

export const data = new SlashCommandBuilder()
  .setName("userinfo-private")
  .setDescription("Envoie en MP une fiche detaillee d'un membre au staff.")
  .setDefaultMemberPermissions(PermissionFlagsBits.ModerateMembers)
  .addUserOption((option) =>
    option.setName("user").setDescription("Membre cible").setRequired(true),
  );

export async function execute(interaction) {
  const guild = interaction.guild;

  if (!guild) {
    await interaction.reply({
      content: "Cette commande doit etre utilisee dans un serveur.",
      ephemeral: true,
    });
    return;
  }

  const targetUser = interaction.options.getUser("user", true);
  const member = await guild.members.fetch(targetUser.id).catch(() => null);
  const warningEntries = warnings.get(`${guild.id}:${targetUser.id}`) ?? [];
  const xpEntry = xp.get(`${guild.id}:${targetUser.id}`) ?? { totalXp: 0, messageCount: 0, level: 0 };

  const embed = new EmbedBuilder()
    .setTitle(`Fiche privee • ${targetUser.tag}`)
    .addFields(
      {
        name: "Identite",
        value: [
          `ID: ${targetUser.id}`,
          `Compte cree: <t:${Math.floor(targetUser.createdTimestamp / 1000)}:F>`,
          `A rejoint: ${
            member?.joinedAt ? `<t:${Math.floor(member.joinedAt.getTime() / 1000)}:F>` : "Inconnu"
          }`,
        ].join("\n"),
        inline: false,
      },
      {
        name: "Moderation",
        value: [
          `Warnings: ${warningEntries.length}`,
          `Derniers warns: ${
            warningEntries.length
              ? warningEntries
                .slice(-3)
                .reverse()
                .map((entry) => `${entry.id} — ${entry.reason}`)
                .join("\n")
              : "Aucun"
          }`,
        ].join("\n"),
        inline: false,
      },
      {
        name: "Activite",
        value: [
          `XP total: ${xpEntry.totalXp ?? 0}`,
          `Niveau: ${xpEntry.level ?? 0}`,
          `Messages comptabilises: ${xpEntry.messageCount ?? 0}`,
        ].join("\n"),
        inline: false,
      },
      {
        name: "Roles",
        value:
          member?.roles.cache
            .filter((role) => role.id !== guild.id)
            .map((role) => role.toString())
            .join(", ") || "Aucun role detaille disponible.",
        inline: false,
      },
    )
    .setThumbnail(targetUser.displayAvatarURL())
    .setColor(0xfaa61a)
    .setTimestamp();

  try {
    await interaction.user.send({
      embeds: [embed],
    });

    await interaction.reply({
      content: `La fiche detaillee de ${targetUser.tag} a ete envoyee en MP.`,
      ephemeral: true,
    });
  } catch {
    await interaction.reply({
      content: "Impossible d'envoyer la fiche en MP. Verifiez que vos messages prives sont ouverts.",
      ephemeral: true,
    });
  }
}
