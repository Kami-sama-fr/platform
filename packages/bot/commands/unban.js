import { EmbedBuilder, PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import { addAudit } from "../utils/store.js";

export const data = new SlashCommandBuilder()
  .setName("unban")
  .setDescription("Debannit un utilisateur du serveur.")
  .setDMPermission(false)
  .setDefaultMemberPermissions(PermissionFlagsBits.BanMembers)
  .addStringOption((o) =>
    o.setName("user")
      .setDescription("Identifiant Discord ou nom d'utilisateur")
      .setRequired(true)
  )
  .addStringOption((o) =>
    o.setName("reason").setDescription("Raison du debannissement")
  );

export async function execute(interaction) {
  if (!interaction.inCachedGuild()) {
    await interaction.reply({ content: "Cette commande doit etre utilisee dans un serveur.", ephemeral: true });
    return;
  }

  const input = interaction.options.getString("user", true).trim();
  const reason = interaction.options.getString("reason") ?? undefined;

  await interaction.deferReply({ ephemeral: true });

  const banList = await interaction.guild.bans.fetch().catch(() => null);
  if (!banList) {
    await interaction.editReply({ content: "Impossible de recuperer la liste des bannis." });
    return;
  }

  const target = input.match(/^\d{17,20}$/)
    ? banList.get(input)
    : banList.find((b) => b.user.username.toLowerCase() === input.toLowerCase() || b.user.tag.toLowerCase() === input.toLowerCase());

  if (!target) {
    await interaction.editReply({
      content: "Aucun utilisateur banni trouve avec cet identifiant ou ce nom.",
    });
    return;
  }

  await interaction.guild.members.unban(target.user.id, reason ? `${interaction.user.tag} — ${reason}` : interaction.user.tag);

  addAudit({
    guildId: interaction.guildId,
    action: "unban",
    actor: interaction.user.tag,
    target: target.user.tag,
    targetId: target.user.id,
    reason: reason ?? "Aucune",
  });

  const embed = new EmbedBuilder()
    .setColor(0x57f287)
    .setTitle("Utilisateur debanni")
    .setDescription(`${target.user.tag} a ete debanni.`)
    .addFields(
      { name: "Utilisateur", value: `${target.user} (\`${target.user.id}\`)`, inline: true },
      { name: "Raison", value: reason ?? "Aucune", inline: true }
    )
    .setFooter({ text: `Par ${interaction.user.tag}` })
    .setTimestamp();

  await interaction.editReply({ embeds: [embed] });
}
