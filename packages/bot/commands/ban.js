import { EmbedBuilder, PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import { addAudit } from "../utils/store.js";

export const data = new SlashCommandBuilder()
  .setName("ban")
  .setDescription("Bannit un membre du serveur.")
  .setDMPermission(false)
  .setDefaultMemberPermissions(PermissionFlagsBits.BanMembers)
  .addUserOption((o) =>
    o.setName("user").setDescription("Membre a bannir").setRequired(true)
  )
  .addStringOption((o) =>
    o.setName("reason").setDescription("Raison du bannissement").setRequired(true)
  )
  .addIntegerOption((o) =>
    o.setName("delete_days")
      .setDescription("Jours de messages a supprimer (0-7)")
      .setMinValue(0)
      .setMaxValue(7)
  );

export async function execute(interaction) {
  if (!interaction.inCachedGuild()) {
    await interaction.reply({ content: "Cette commande doit etre utilisee dans un serveur.", ephemeral: true });
    return;
  }

  const user = interaction.options.getUser("user", true);
  const reason = interaction.options.getString("reason", true);
  const deleteDays = interaction.options.getInteger("delete_days") ?? 0;
  const member = await interaction.guild.members.fetch(user.id).catch(() => null);

  if (!member) {
    await interaction.reply({ content: "Ce membre n'est pas sur le serveur.", ephemeral: true });
    return;
  }

  if (user.id === interaction.client.user.id) {
    await interaction.reply({ content: "Je ne peux pas me bannir moi-meme.", ephemeral: true });
    return;
  }

  if (user.id === interaction.guild.ownerId) {
    await interaction.reply({ content: "Impossible de bannir le proprietaire du serveur.", ephemeral: true });
    return;
  }

  if (!member.bannable) {
    await interaction.reply({ content: "Ce membre ne peut pas etre banni par le bot (verifiez la hierarchie des roles).", ephemeral: true });
    return;
  }

  if (interaction.member.roles.highest.position <= member.roles.highest.position && interaction.guild.ownerId !== interaction.user.id) {
    await interaction.reply({ content: "Vous ne pouvez pas bannir un membre avec un role egal ou superieur au votre.", ephemeral: true });
    return;
  }

  await interaction.deferReply({ ephemeral: true });

  try {
    await member.send({
      embeds: [
        new EmbedBuilder()
          .setColor(0xed4245)
          .setTitle("Vous avez ete banni")
          .setDescription(`Vous avez ete banni de **${interaction.guild.name}**.`)
          .addFields({ name: "Raison", value: reason })
          .setFooter({ text: `Bannissement • ${interaction.user.tag}` })
          .setTimestamp(),
      ],
    }).catch(() => {});
  } catch {}

  await interaction.guild.members.ban(user, {
    reason: `${interaction.user.tag} — ${reason}`,
    deleteMessageSeconds: deleteDays * 86400,
  });

  addAudit({
    guildId: interaction.guildId,
    action: "ban",
    actor: interaction.user.tag,
    target: user.tag,
    targetId: user.id,
    reason,
    deleteDays,
  });

  const embed = new EmbedBuilder()
    .setColor(0xed4245)
    .setTitle("Membre banni")
    .setDescription(`${user.tag} a ete banni.`)
    .addFields(
      { name: "Utilisateur", value: `${user} (\`${user.id}\`)`, inline: true },
      { name: "Raison", value: reason, inline: true },
      { name: "Messages supprimes", value: deleteDays > 0 ? `${deleteDays} jour(s)` : "Aucun", inline: true }
    )
    .setFooter({ text: `Par ${interaction.user.tag}` })
    .setTimestamp();

  await interaction.editReply({ embeds: [embed] });
}
