import { SlashCommandBuilder } from "discord.js";
import { getGuildLeaderboard } from "../utils/xp.js";

export const data = new SlashCommandBuilder()
  .setName("leaderboard")
  .setDescription("Affiche le classement XP du serveur.")
  .addIntegerOption((option) =>
    option
      .setName("limit")
      .setDescription("Nombre de membres à afficher")
      .setMinValue(3)
      .setMaxValue(20)
      .setRequired(false),
  );

export async function execute(interaction) {
  if (!interaction.guild) {
    await interaction.reply({
      content: "Cette commande doit être utilisée dans un serveur.",
      ephemeral: true,
    });
    return;
  }

  const limit = interaction.options.getInteger("limit") ?? 10;
  const leaderboard = getGuildLeaderboard(interaction.guild.id, limit);

  if (leaderboard.length === 0) {
    await interaction.reply({
      content: "Aucune donnée XP n'est encore disponible sur ce serveur.",
      ephemeral: true,
    });
    return;
  }

  const lines = await Promise.all(
    leaderboard.map(async (entry) => {
      const user = await interaction.client.users.fetch(entry.userId).catch(() => null);
      return `#${entry.rank} - ${user ? user.tag : entry.userId} • Niveau ${entry.level} • ${entry.totalXp} XP`;
    }),
  );

  await interaction.reply({
    content: `Classement XP de **${interaction.guild.name}**:\n\n${lines.join("\n")}`,
    ephemeral: false,
  });
}
