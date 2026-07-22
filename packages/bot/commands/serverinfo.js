import { SlashCommandBuilder } from "discord.js";

export const data = new SlashCommandBuilder()
  .setName("serverinfo")
  .setDescription("Affiche les informations du serveur Discord.");

export async function execute(interaction) {
  const guild = interaction.guild;

  if (!guild) {
    await interaction.reply({
      content: "Cette commande doit être utilisée dans un serveur.",
      ephemeral: true,
    });
    return;
  }

  await interaction.reply({
    content:
      `Serveur: ${guild.name}\n` +
      `ID: ${guild.id}\n` +
      `Membres: ${guild.memberCount}\n` +
      `Créé le: <t:${Math.floor(guild.createdTimestamp / 1000)}:F>`,
    ephemeral: true,
  });
}
