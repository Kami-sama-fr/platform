import { ChannelType, EmbedBuilder, PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import { addAudit } from "../utils/store.js";

export const data = new SlashCommandBuilder()
  .setName("slowmode")
  .setDescription("Configure le mode lent d'un salon.")
  .setDMPermission(false)
  .setDefaultMemberPermissions(PermissionFlagsBits.ManageChannels)
  .addIntegerOption((o) =>
    o.setName("seconds")
      .setDescription("Duree en secondes (0 pour desactiver)")
      .setRequired(true)
      .setMinValue(0)
      .setMaxValue(21600)
  )
  .addChannelOption((o) =>
    o.setName("channel").setDescription("Salon cible (defaut: celui-ci)")
  )
  .addStringOption((o) =>
    o.setName("reason").setDescription("Raison de la modification")
  );

export async function execute(interaction) {
  if (!interaction.inCachedGuild()) {
    await interaction.reply({ content: "Cette commande doit etre utilisee dans un serveur.", ephemeral: true });
    return;
  }

  const seconds = interaction.options.getInteger("seconds", true);
  const channel = interaction.options.getChannel("channel") ?? interaction.channel;
  const reason = interaction.options.getString("reason") ?? undefined;

  if (!channel || !channel.isTextBased() || channel.isDMBased()) {
    await interaction.reply({ content: "Salon invalide ou incompatible.", ephemeral: true });
    return;
  }

  if (![ChannelType.GuildText, ChannelType.GuildAnnouncement, ChannelType.PublicThread, ChannelType.PrivateThread].includes(channel.type)) {
    await interaction.reply({ content: "Ce type de salon ne supporte pas le mode lent.", ephemeral: true });
    return;
  }

  const botMember = interaction.guild.members.me;
  const perms = channel.permissionsFor(botMember);
  if (!perms?.has(PermissionFlagsBits.ManageChannels)) {
    await interaction.reply({ content: `Le bot n'a pas la permission de gerer ce salon.`, ephemeral: true });
    return;
  }

  if (channel.rateLimitPerUser === seconds) {
    await interaction.reply({ content: `Le mode lent est deja regle sur ${seconds}s dans ${channel}.`, ephemeral: true });
    return;
  }

  await channel.setRateLimitPerUser(seconds, reason ?? `Par ${interaction.user.tag}`);

  addAudit({
    guildId: interaction.guildId,
    action: "slowmode",
    actor: interaction.user.tag,
    target: channel.id,
    reason: reason ?? `${seconds}s`,
  });

  const embed = new EmbedBuilder()
    .setColor(seconds > 0 ? 0xfee75c : 0x57f287)
    .setTitle("Mode lent configure")
    .setDescription(seconds > 0
      ? `Le mode lent est desormais de **${seconds}s** dans ${channel}.`
      : `Le mode lent a ete desactive dans ${channel}.`
    )
    .addFields(
      { name: "Salon", value: `${channel}`, inline: true },
      { name: "Duree", value: seconds > 0 ? `${seconds} seconde(s)` : "Desactive", inline: true },
    )
    .setFooter({ text: `Par ${interaction.user.tag}` })
    .setTimestamp();

  await interaction.reply({ embeds: [embed], ephemeral: true });
}
