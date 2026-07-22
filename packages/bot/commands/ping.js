import { SlashCommandBuilder } from "discord.js";

export const data = new SlashCommandBuilder()
  .setName("ping")
  .setDescription("Affiche la latence du bot.");

export async function execute(interaction) {
  const gatewayPing = Math.round(interaction.client.ws.ping);
  const apiLatency = Date.now() - interaction.createdTimestamp;

  await interaction.reply({
    content: `Pong! Gateway: ${gatewayPing}ms | API: ${apiLatency}ms`,
    ephemeral: true,
  });
}
