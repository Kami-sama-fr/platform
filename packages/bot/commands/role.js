import { PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import { addAudit } from "../utils/store.js";

export const data = new SlashCommandBuilder()
  .setName("role")
  .setDescription("Ajoute ou retire un rôle du serveur à un membre.")
  .setDefaultMemberPermissions(PermissionFlagsBits.ManageRoles)
  .addSubcommand((subcommand) =>
    subcommand
      .setName("add")
      .setDescription("Ajouter un rôle à un membre")
      .addUserOption((option) =>
        option.setName("user").setDescription("Membre").setRequired(true),
      )
      .addRoleOption((option) =>
        option.setName("role").setDescription("Rôle du serveur").setRequired(true),
      ),
  )
  .addSubcommand((subcommand) =>
    subcommand
      .setName("remove")
      .setDescription("Retirer un rôle d'un membre")
      .addUserOption((option) =>
        option.setName("user").setDescription("Membre").setRequired(true),
      )
      .addRoleOption((option) =>
        option.setName("role").setDescription("Rôle du serveur").setRequired(true),
      ),
  );

export async function execute(interaction) {
  const action = interaction.options.getSubcommand();
  const member = interaction.options.getMember("user");
  const role = interaction.options.getRole("role");

  if (!member || !role) {
    await interaction.reply({
      content: "Membre ou rôle introuvable.",
      ephemeral: true,
    });
    return;
  }

  if (role.id === interaction.guild.id) {
    await interaction.reply({
      content: "Le rôle @everyone ne peut pas être attribué avec cette commande.",
      ephemeral: true,
    });
    return;
  }

  if (role.managed) {
    await interaction.reply({
      content: "Ce rôle est géré par une intégration et ne peut pas être modifié ici.",
      ephemeral: true,
    });
    return;
  }

  const botMember = interaction.guild.members.me;

  if (!botMember) {
    await interaction.reply({
      content: "Impossible de vérifier les permissions du bot sur ce serveur.",
      ephemeral: true,
    });
    return;
  }

  if (role.position >= botMember.roles.highest.position) {
    await interaction.reply({
      content: "Le bot ne peut pas gérer ce rôle car il est au-dessus ou au même niveau que son rôle le plus élevé.",
      ephemeral: true,
    });
    return;
  }

  if (role.position >= interaction.member.roles.highest.position && interaction.guild.ownerId !== interaction.user.id) {
    await interaction.reply({
      content: "Vous ne pouvez pas gérer un rôle situé au-dessus ou au même niveau que votre rôle le plus élevé.",
      ephemeral: true,
    });
    return;
  }

  if (action === "add") {
    if (member.roles.cache.has(role.id)) {
      await interaction.reply({
        content: `${member.user.tag} possède déjà le rôle ${role.name}.`,
        ephemeral: true,
      });
      return;
    }

    await member.roles.add(role);
    addAudit({
      guildId: interaction.guildId,
      action: "role_add",
      actor: interaction.user.tag,
      target: member.user.tag,
      reason: role.name,
    });
    await interaction.reply({
      content: `Le rôle ${role.name} a été ajouté à ${member.user.tag}.`,
      ephemeral: true,
    });
    return;
  }

  if (!member.roles.cache.has(role.id)) {
    await interaction.reply({
      content: `${member.user.tag} ne possède pas le rôle ${role.name}.`,
      ephemeral: true,
    });
    return;
  }

  await member.roles.remove(role);
  addAudit({
    guildId: interaction.guildId,
    action: "role_remove",
    actor: interaction.user.tag,
    target: member.user.tag,
    reason: role.name,
  });
  await interaction.reply({
    content: `Le rôle ${role.name} a été retiré de ${member.user.tag}.`,
    ephemeral: true,
  });
}
