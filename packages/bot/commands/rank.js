import { SlashCommandBuilder } from "discord.js";
import { getGuildRank, getLevelFromXp, getMemberXp } from "../utils/xp.js";

export const data = new SlashCommandBuilder()
  .setName("rank")
  .setDescription("Affiche le niveau et l'XP d'un membre.")
  .addUserOption((option) =>
    option.setName("user").setDescription("Membre cible").setRequired(false),
  );

export async function execute(interaction) {
  if (!interaction.guild) {
    await interaction.reply({
      content: "Cette commande doit être utilisée dans un serveur.",
      ephemeral: true,
    });
    return;
  }

  const targetUser = interaction.options.getUser("user") ?? interaction.user;
  const xpEntry = getMemberXp(interaction.guild.id, targetUser.id);
  const levelState = getLevelFromXp(xpEntry.totalXp);
  const rank = getGuildRank(interaction.guild.id, targetUser.id);

  await interaction.reply({
    content:
      `Classement de ${targetUser}:\n` +
      `Niveau: ${levelState.level}\n` +
      `XP total: ${xpEntry.totalXp}\n` +
      `Progression: ${levelState.currentLevelXp}/${levelState.nextLevelXp}\n` +
      `Rang serveur: ${rank ? `#${rank}` : "Non classé"}\n` +
      `Messages comptabilisés: ${xpEntry.messageCount}`,
    ephemeral: false,
  });
}
