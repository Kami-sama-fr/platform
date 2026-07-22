import { EmbedBuilder, PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import { addAudit, warnings } from "../utils/store.js";

export const data = new SlashCommandBuilder()
  .setName("warn")
  .setDescription("Ajoute un avertissement a un membre.")
  .setDMPermission(false)
  .setDefaultMemberPermissions(PermissionFlagsBits.ModerateMembers)
  .addUserOption((o) =>
    o.setName("user").setDescription("Membre a avertir").setRequired(true)
  )
  .addStringOption((o) =>
    o.setName("reason").setDescription("Raison de l'avertissement").setRequired(true)
  );

export async function execute(interaction) {
  if (!interaction.inCachedGuild()) {
    await interaction.reply({ content: "Cette commande doit etre utilisee dans un serveur.", ephemeral: true });
    return;
  }

  const user = interaction.options.getUser("user", true);
  const reason = interaction.options.getString("reason", true);
  const member = await interaction.guild.members.fetch(user.id).catch(() => null);

  if (!member) {
    await interaction.reply({ content: "Ce membre n'est pas sur le serveur.", ephemeral: true });
    return;
  }

  if (user.id === interaction.client.user.id) {
    await interaction.reply({ content: "Je ne peux pas m'avertir moi-meme.", ephemeral: true });
    return;
  }

  if (user.id === interaction.guild.ownerId) {
    await interaction.reply({ content: "Impossible d'avertir le proprietaire du serveur.", ephemeral: true });
    return;
  }

  if (interaction.member.roles.highest.position <= member.roles.highest.position && interaction.guild.ownerId !== interaction.user.id) {
    await interaction.reply({ content: "Vous ne pouvez pas avertir un membre avec un role egal ou superieur au votre.", ephemeral: true });
    return;
  }

  const key = `${interaction.guildId}:${user.id}`;
  const list = warnings.get(key) ?? [];
  const item = {
    id: crypto.randomUUID().slice(0, 8),
    reason,
    moderator: interaction.user.tag,
    at: new Date(),
  };
  list.push(item);
  warnings.set(key, list);

  addAudit({
    guildId: interaction.guildId,
    action: "warn",
    actor: interaction.user.tag,
    target: user.tag,
    targetId: user.id,
    reason,
    warnId: item.id,
    total: list.length,
  });

  try {
    await member.send({
      embeds: [
        new EmbedBuilder()
          .setColor(0xfee75c)
          .setTitle("Avertissement recu")
          .setDescription(`Vous avez recu un avertissement sur **${interaction.guild.name}**.`)
          .addFields(
            { name: "Avertissement", value: `\`${item.id}\``, inline: true },
            { name: "Raison", value: reason, inline: true }
          )
          .setFooter({ text: interaction.user.tag })
          .setTimestamp(),
      ],
    }).catch(() => {});
  } catch {}

  const embed = new EmbedBuilder()
    .setColor(0xfee75c)
    .setTitle("Avertissement ajoute")
    .setDescription(`${user} a recu un avertissement.`)
    .addFields(
      { name: "Utilisateur", value: `${user.tag} (\`${user.id}\`)`, inline: true },
      { name: "Avertissement", value: `\`${item.id}\``, inline: true },
      { name: "Raison", value: reason, inline: false },
      { name: "Total des avertissements", value: String(list.length), inline: true }
    )
    .setFooter({ text: `Par ${interaction.user.tag}` })
    .setTimestamp();

  await interaction.reply({ embeds: [embed], ephemeral: true });
}
