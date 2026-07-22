import { SlashCommandBuilder } from "discord.js";

export const data = new SlashCommandBuilder()
  .setName("about")
  .setDescription("Présente Sky Genesis Enterprise.");

export async function execute(interaction) {
  await interaction.reply({
    content:
      "Sky Genesis Enterprise (SGE) est l'écosystème officiel autour des services Aether. Ce serveur Discord sert de hub communautaire pour les annonces, le support et la collaboration entre membres, développeurs et partenaires.",
    ephemeral: false,
  });
}
