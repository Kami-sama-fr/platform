import { PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import { addAudit } from "../utils/store.js";
import { addMemberXp, getLevelFromXp, getMemberXp, resetMemberXp, setMemberXp } from "../utils/xp.js";

export const data = new SlashCommandBuilder()
  .setName("xp")
  .setDescription("Administre l'XP des membres.")
  .setDefaultMemberPermissions(PermissionFlagsBits.ManageGuild)
  .addSubcommand((subcommand) =>
    subcommand
      .setName("add")
      .setDescription("Ajoute de l'XP à un membre.")
      .addUserOption((option) => option.setName("user").setDescription("Membre cible").setRequired(true))
      .addIntegerOption((option) =>
        option.setName("amount").setDescription("Quantité d'XP").setMinValue(1).setRequired(true),
      ),
  )
  .addSubcommand((subcommand) =>
    subcommand
      .setName("remove")
      .setDescription("Retire de l'XP à un membre.")
      .addUserOption((option) => option.setName("user").setDescription("Membre cible").setRequired(true))
      .addIntegerOption((option) =>
        option.setName("amount").setDescription("Quantité d'XP").setMinValue(1).setRequired(true),
      ),
  )
  .addSubcommand((subcommand) =>
    subcommand
      .setName("set")
      .setDescription("Définit l'XP exacte d'un membre.")
      .addUserOption((option) => option.setName("user").setDescription("Membre cible").setRequired(true))
      .addIntegerOption((option) =>
        option.setName("amount").setDescription("XP totale").setMinValue(0).setRequired(true),
      ),
  )
  .addSubcommand((subcommand) =>
    subcommand
      .setName("reset")
      .setDescription("Réinitialise le profil XP d'un membre.")
      .addUserOption((option) => option.setName("user").setDescription("Membre cible").setRequired(true)),
  );

export async function execute(interaction) {
  if (!interaction.guild) {
    await interaction.reply({
      content: "Cette commande doit être utilisée dans un serveur.",
      ephemeral: true,
    });
    return;
  }

  const subcommand = interaction.options.getSubcommand();
  const targetUser = interaction.options.getUser("user", true);
  const amount = interaction.options.getInteger("amount");

  if (subcommand === "reset") {
    resetMemberXp(interaction.guild.id, targetUser.id);
    addAudit({
      guildId: interaction.guild.id,
      action: "xp_reset",
      actor: interaction.user.tag,
      target: targetUser.tag,
    });

    await interaction.reply({
      content: `Le profil XP de ${targetUser.tag} a été réinitialisé.`,
      ephemeral: true,
    });
    return;
  }

  let result;

  if (subcommand === "add") {
    result = addMemberXp(interaction.guild.id, targetUser.id, amount ?? 0);
  } else if (subcommand === "remove") {
    result = addMemberXp(interaction.guild.id, targetUser.id, -(amount ?? 0));
  } else {
    const current = getMemberXp(interaction.guild.id, targetUser.id);
    result = setMemberXp(interaction.guild.id, targetUser.id, amount ?? 0, current.messageCount);
  }

  const levelState = getLevelFromXp(result.totalXp);

  addAudit({
    guildId: interaction.guild.id,
    action: `xp_${subcommand}`,
    actor: interaction.user.tag,
    target: targetUser.tag,
    reason: `${amount ?? 0} XP`,
  });

  await interaction.reply({
    content:
      `${targetUser.tag} a maintenant ${result.totalXp} XP.\n` +
      `Niveau: ${levelState.level}\n` +
      `Progression: ${levelState.currentLevelXp}/${levelState.nextLevelXp}`,
    ephemeral: true,
  });
}
