import { EmbedBuilder, PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import { addAudit } from "../utils/store.js";

export const data = new SlashCommandBuilder()
  .setName("nickname")
  .setDescription("Modifie le surnom d'un membre.")
  .setDMPermission(false)
  .setDefaultMemberPermissions(PermissionFlagsBits.ManageNicknames)
  .addSubcommand((s) =>
    s
      .setName("set")
      .setDescription("Definit un surnom")
      .addUserOption((o) => o.setName("user").setDescription("Membre").setRequired(true))
      .addStringOption((o) =>
        o.setName("name").setDescription("Nouveau surnom").setRequired(true).setMaxLength(32)
      )
      .addStringOption((o) => o.setName("reason").setDescription("Raison du changement"))
  )
  .addSubcommand((s) =>
    s
      .setName("reset")
      .setDescription("Reinitialise le surnom")
      .addUserOption((o) => o.setName("user").setDescription("Membre").setRequired(true))
      .addStringOption((o) => o.setName("reason").setDescription("Raison de la reinitialisation"))
  );

export async function execute(interaction) {
  if (!interaction.inCachedGuild()) {
    await interaction.reply({ content: "Cette commande doit etre utilisee dans un serveur.", ephemeral: true });
    return;
  }

  const member = interaction.options.getMember("user");
  if (!member) {
    await interaction.reply({ content: "Membre introuvable sur le serveur.", ephemeral: true });
    return;
  }

  if (member.id === interaction.guild.ownerId && interaction.user.id !== interaction.guild.ownerId) {
    await interaction.reply({ content: "Vous ne pouvez pas modifier le surnom du proprietaire du serveur.", ephemeral: true });
    return;
  }

  if (interaction.member.roles.highest.position <= member.roles.highest.position && interaction.guild.ownerId !== interaction.user.id) {
    await interaction.reply({ content: "Vous ne pouvez pas modifier le surnom d'un membre avec un role egal ou superieur au votre.", ephemeral: true });
    return;
  }

  const sub = interaction.options.getSubcommand();
  const reason = interaction.options.getString("reason") ?? undefined;
  const previousNickname = member.nickname ?? member.displayName;
  const name = sub === "set" ? interaction.options.getString("name", true) : null;

  if (name !== null && name === previousNickname) {
    await interaction.reply({ content: `Le surnom est deja defini sur \`${name}\`.`, ephemeral: true });
    return;
  }

  const botMember = interaction.guild.members.me;
  if (botMember.roles.highest.position <= member.roles.highest.position) {
    await interaction.reply({ content: "Le bot ne peut pas modifier ce membre (role trop eleve).", ephemeral: true });
    return;
  }

  const displayName = name ?? member.user.username;
  await member.setNickname(name, reason ?? `${interaction.user.tag} — ${sub}`);

  addAudit({
    guildId: interaction.guildId,
    action: sub === "set" ? "nickname_set" : "nickname_reset",
    actor: interaction.user.tag,
    target: member.user.tag,
    targetId: member.id,
    reason: reason ?? "Aucune",
    previous: previousNickname,
    new: displayName,
  });

  const embed = new EmbedBuilder()
    .setColor(0x5865f2)
    .setTitle(sub === "set" ? "Surnom modifie" : "Surnom reinitialise")
    .setDescription(sub === "set"
      ? `Le surnom de ${member.user} a ete defini sur \`${displayName}\`.`
      : `Le surnom de ${member.user} a ete reinitialise.`
    )
    .addFields(
      { name: "Membre", value: `${member.user} (\`${member.id}\`)`, inline: true },
      { name: "Ancien surnom", value: `\`${previousNickname}\``, inline: true },
      { name: "Nouveau surnom", value: `\`${displayName}\``, inline: true },
    )
    .setFooter({ text: `Par ${interaction.user.tag}` })
    .setTimestamp();

  await interaction.reply({ embeds: [embed], ephemeral: true });
}
