import { EmbedBuilder, PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import { auditLog, incidents, tickets, warnings } from "../utils/store.js";

export const data = new SlashCommandBuilder()
  .setName("admin-report")
  .setDescription("Envoie en MP un recap prive de moderation et d'activite.")
  .setDefaultMemberPermissions(PermissionFlagsBits.Administrator);

export async function execute(interaction) {
  const guild = interaction.guild;

  if (!guild) {
    await interaction.reply({
      content: "Cette commande doit etre utilisee dans un serveur.",
      ephemeral: true,
    });
    return;
  }

  const guildWarnings = [...warnings.entries()]
    .filter(([key]) => key.startsWith(`${guild.id}:`))
    .flatMap(([, entries]) => entries);

  const guildTickets = [...tickets.entries()]
    .filter(([, ticket]) => ticket.guildId === guild.id);

  const guildIncidents = [...incidents.values()].filter(
    (incident) => !incident.guildId || incident.guildId === guild.id,
  );

  const recentAudit = auditLog
    .filter((entry) => entry.guildId === guild.id)
    .slice(0, 5)
    .map(
      (entry) =>
        `• ${entry.action} — ${entry.actor}${entry.target ? ` -> ${entry.target}` : ""} — <t:${Math.floor(entry.at.getTime() / 1000)}:R>`,
    )
    .join("\n");

  const recentJoins = guild.members.cache
    .filter((member) => !member.user.bot && member.joinedAt)
    .sort((left, right) => (right.joinedAt?.getTime() ?? 0) - (left.joinedAt?.getTime() ?? 0))
    .first(5)
    .map((member) => `${member.user.tag} — <t:${Math.floor(member.joinedAt.getTime() / 1000)}:R>`)
    .join("\n");

  const embed = new EmbedBuilder()
    .setTitle(`Rapport admin prive • ${guild.name}`)
    .setDescription("Resume rapide des signaux utiles pour la moderation et le support.")
    .addFields(
      {
        name: "Vue d'ensemble",
        value: [
          `Membres: ${guild.memberCount}`,
          `Tickets ouverts: ${guildTickets.length}`,
          `Warnings stockes: ${guildWarnings.length}`,
          `Incidents suivis: ${guildIncidents.length}`,
        ].join("\n"),
        inline: false,
      },
      {
        name: "Dernieres actions staff",
        value: recentAudit || "Aucune action recente enregistre.",
        inline: false,
      },
      {
        name: "Dernieres arrivees",
        value: recentJoins || "Aucune arrivee recente disponible dans le cache.",
        inline: false,
      },
      {
        name: "Tickets ouverts",
        value:
          guildTickets
            .slice(0, 5)
            .map(
              ([channelId, ticket]) =>
                `• <#${channelId}> — owner <@${ticket.owner}>${ticket.topic ? ` — ${ticket.topic}` : ""}`,
            )
            .join("\n") || "Aucun ticket ouvert.",
        inline: false,
      },
    )
    .setColor(0x5865f2)
    .setTimestamp();

  try {
    await interaction.user.send({
      embeds: [embed],
    });

    await interaction.reply({
      content: "Le rapport admin a ete envoye en message prive.",
      ephemeral: true,
    });
  } catch {
    await interaction.reply({
      content: "Impossible d'envoyer le rapport en MP. Verifiez que vos messages prives sont ouverts.",
      ephemeral: true,
    });
  }
}
