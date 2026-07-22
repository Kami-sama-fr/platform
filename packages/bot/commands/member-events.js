import { PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import {
  addAudit,
  areMemberEventsEnabled,
  setMemberEventsEnabled,
  updateGoodbyeSettings,
  updateWelcomeSettings,
} from "../utils/store.js";

export const data = new SlashCommandBuilder()
  .setName("member-events")
  .setDescription("Active ou desactive l'accueil automatique des membres.")
  .setDefaultMemberPermissions(PermissionFlagsBits.ManageGuild)
  .addSubcommand((subcommand) =>
    subcommand
      .setName("enable")
      .setDescription("Active les messages d'arrivee/depart et le role automatique.")
  )
  .addSubcommand((subcommand) =>
    subcommand
      .setName("disable")
      .setDescription("Desactive les messages d'arrivee/depart et le role automatique.")
  )
  .addSubcommand((subcommand) =>
    subcommand.setName("status").setDescription("Affiche l'etat actuel de la fonction.")
  );

export async function execute(interaction) {
  if (!interaction.guild) {
    await interaction.reply({
      content: "Cette commande doit etre utilisee dans un serveur.",
      ephemeral: true,
    });
    return;
  }

  const subcommand = interaction.options.getSubcommand();

  if (subcommand === "status") {
    const enabled = areMemberEventsEnabled(interaction.guild.id);

    await interaction.reply({
      content: `L'accueil automatique des membres est actuellement **${enabled ? "active" : "desactive"}**.`,
      ephemeral: true,
    });
    return;
  }

  const enabled = subcommand === "enable";
  setMemberEventsEnabled(interaction.guild.id, enabled);
  updateWelcomeSettings(interaction.guild.id, { enabled });
  updateGoodbyeSettings(interaction.guild.id, { enabled });
  addAudit({
    guildId: interaction.guild.id,
    action: enabled ? "member_events_enable" : "member_events_disable",
    actor: interaction.user.tag,
  });

  await interaction.reply({
    content: enabled
      ? "L'accueil automatique est active. Les arrivees, departs et attributions du role membre seront traites."
      : "L'accueil automatique est desactive. Aucun message ni role automatique ne sera applique.",
    ephemeral: true,
  });
}
