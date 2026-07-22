import { SlashCommandBuilder, PermissionFlagsBits } from "discord.js";

const BULK_DELETE_MAX_AGE_MS = 14 * 24 * 60 * 60 * 1000;

export const data = new SlashCommandBuilder()
  .setName("clear")
  .setDescription("Supprime des messages récents.")
  .setDMPermission(false)
  .setDefaultMemberPermissions(PermissionFlagsBits.ManageMessages)
  .addIntegerOption((option) =>
    option
      .setName("amount")
      .setDescription("Nombre de messages")
      .setRequired(true)
      .setMinValue(1)
      .setMaxValue(100),
  );

export async function execute(interaction) {
  if (!interaction.channel?.bulkDelete || !interaction.channel.messages?.fetch) {
    return interaction.reply({ content: "Salon incompatible.", ephemeral: true });
  }

  const requestedAmount = interaction.options.getInteger("amount", true);
  const fetchedMessages = await interaction.channel.messages.fetch({ limit: requestedAmount });
  const newestMessages = fetchedMessages.filter(
    (message) => Date.now() - message.createdTimestamp < BULK_DELETE_MAX_AGE_MS,
  );

  if (newestMessages.size === 0) {
    return interaction.reply({
      content: "Aucun message recent supprimable n'a ete trouve dans ce salon.",
      ephemeral: true,
    });
  }

  const deletedMessages = await interaction.channel.bulkDelete(newestMessages, true);
  const missingCount = Math.max(requestedAmount - deletedMessages.size, 0);

  if (missingCount > 0) {
    return interaction.reply({
      content: `${deletedMessages.size} message(s) supprime(s). ${missingCount} message(s) demande(s) etaient absents ou trop anciens pour une suppression groupée.`,
      ephemeral: true,
    });
  }

  return interaction.reply({
    content: `${deletedMessages.size} message(s) supprime(s).`,
    ephemeral: true,
  });
}
