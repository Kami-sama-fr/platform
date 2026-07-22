import { SlashCommandBuilder } from "discord.js";

export const data = new SlashCommandBuilder()
  .setName("help")
  .setDescription("Affiche la liste des commandes disponibles.");

export async function execute(interaction) {
  const commands = [...interaction.client.commands.values()]
    .map((command) => `• /${command.data.name} - ${command.data.description}`)
    .join("\n");

  await interaction.reply({
    content: `Voici les commandes disponibles:\n\n${commands}`,
    ephemeral: true,
  });
}
