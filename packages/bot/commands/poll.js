import { EmbedBuilder, PermissionFlagsBits, SlashCommandBuilder } from "discord.js";
import { addAudit, getPoll, setPoll } from "../utils/store.js";

const EMOJIS = ["1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"];

function parseChoices(raw) {
  return raw
    .split("|")
    .map((x) => x.trim())
    .filter(Boolean)
    .slice(0, 10);
}

function buildPollEmbed(question, choices, guild, author) {
  return new EmbedBuilder()
    .setColor(0x5865f2)
    .setTitle("📊 " + question)
    .setDescription(choices.map((c, i) => `${EMOJIS[i]} ${c}`).join("\n\n"))
    .addFields(
      { name: "Choix", value: String(choices.length), inline: true },
      { name: "Serveur", value: guild.name, inline: true }
    )
    .setFooter({ text: `Sondage de ${author.tag}` })
    .setTimestamp();
}

function buildResultEmbed(question, choices, results, total, guild, author, ended) {
  const maxVotes = Math.max(...results, 1);
  const lines = choices.map((c, i) => {
    const count = results[i] ?? 0;
    const barLen = Math.round((count / maxVotes) * 12);
    const bar = "▰".repeat(barLen) + "▱".repeat(12 - barLen);
    const pct = total > 0 ? ((count / total) * 100).toFixed(1) : "0.0";
    return `${EMOJIS[i]} **${c}**\n${bar} ${count} voix (${pct}%)`;
  });

  return new EmbedBuilder()
    .setColor(ended ? 0xed4245 : 0x57f287)
    .setTitle(ended ? "📊 Sondage termine" : "📊 Resultats du sondage")
    .setDescription(ended
      ? `**${question}**\n\n${lines.join("\n\n")}`
      : `**${question}**\n\n${lines.join("\n\n")}`
    )
    .addFields(
      { name: "Total des votes", value: String(total), inline: true },
      { name: "Statut", value: ended ? "Termine" : "En cours", inline: true }
    )
    .setFooter({ text: ended ? `Sondage de ${author.tag}` : `Resultats en direct • ${author.tag}` })
    .setTimestamp();
}

export const data = new SlashCommandBuilder()
  .setName("poll")
  .setDescription("Cree et gere des sondages.")
  .setDMPermission(false)
  .addSubcommand((sub) =>
    sub
      .setName("create")
      .setDescription("Cree un nouveau sondage.")
      .addStringOption((o) =>
        o.setName("question").setDescription("Question du sondage").setRequired(true)
      )
      .addStringOption((o) =>
        o.setName("choices")
          .setDescription("Choix separes par | (max 10)")
          .setRequired(true)
      )
      .addIntegerOption((o) =>
        o.setName("duration")
          .setDescription("Duree en minutes (0 = illimite)")
          .setMinValue(0)
          .setMaxValue(1440)
      )
  )
  .addSubcommand((sub) =>
    sub
      .setName("end")
      .setDescription("Termine un sondage et affiche les resultats.")
      .addStringOption((o) =>
        o.setName("message_id")
          .setDescription("ID du message du sondage")
          .setRequired(true)
      )
  )
  .addSubcommand((sub) =>
    sub
      .setName("results")
      .setDescription("Affiche les resultats d'un sondage.")
      .addStringOption((o) =>
        o.setName("message_id")
          .setDescription("ID du message du sondage")
          .setRequired(true)
      )
  );

export async function execute(interaction) {
  if (!interaction.inCachedGuild()) {
    await interaction.reply({ content: "Cette commande doit etre utilisee dans un serveur.", ephemeral: true });
    return;
  }

  const sub = interaction.options.getSubcommand();

  if (sub === "create") {
    await handleCreate(interaction);
  } else if (sub === "end") {
    await handleEnd(interaction);
  } else if (sub === "results") {
    await handleResults(interaction);
  }
}

async function handleCreate(interaction) {
  const question = interaction.options.getString("question", true);
  const raw = interaction.options.getString("choices", true);
  const duration = interaction.options.getInteger("duration") ?? 0;
  const choices = parseChoices(raw);

  if (choices.length < 2) {
    await interaction.reply({ content: "Indique au moins deux choix separes par `|`.", ephemeral: true });
    return;
  }

  const embed = buildPollEmbed(question, choices, interaction.guild, interaction.user);

  const msg = await interaction.reply({
    embeds: [embed],
    fetchReply: true,
  });

  for (let i = 0; i < choices.length; i++) {
    await msg.react(EMOJIS[i]);
  }

  const key = `${interaction.guildId}:${msg.id}`;
  const pollData = {
    guildId: interaction.guildId,
    channelId: interaction.channelId,
    messageId: msg.id,
    question,
    choices,
    authorId: interaction.user.id,
    authorTag: interaction.user.tag,
    createdAt: Date.now(),
    expiresAt: duration > 0 ? Date.now() + duration * 60000 : null,
    ended: false,
  };

  setPoll(key, pollData);

  addAudit({
    guildId: interaction.guildId,
    action: "poll_create",
    actor: interaction.user.tag,
    target: msg.id,
    reason: question,
  });

  if (duration > 0) {
    setTimeout(async () => {
      const current = getPoll(key);
      if (current && !current.ended) {
        current.ended = true;
        setPoll(key, current);
        try {
          const channel = await interaction.guild.channels.fetch(current.channelId);
          if (channel?.isTextBased()) {
            const message = await channel.messages.fetch(current.messageId).catch(() => null);
            if (message) {
              const results = await tallyVotes(message, current.choices);
              const total = results.reduce((a, b) => a + b, 0);
              const resultEmbed = buildResultEmbed(
                current.question, current.choices, results, total,
                interaction.guild, { tag: current.authorTag }, true
              );
              await message.edit({ embeds: [resultEmbed] });
            }
          }
        } catch {}
      }
    }, duration * 60000);
  }
}

async function handleEnd(interaction) {
  const messageId = interaction.options.getString("message_id", true);
  const key = `${interaction.guildId}:${messageId}`;
  const poll = getPoll(key);

  if (!interaction.memberPermissions?.has(PermissionFlagsBits.ManageMessages)) {
    await interaction.reply({ content: "Vous n'avez pas la permission de gerer les messages.", ephemeral: true });
    return;
  }

  if (!poll) {
    await interaction.reply({ content: "Aucun sondage actif trouve avec cet ID.", ephemeral: true });
    return;
  }

  if (poll.ended) {
    await interaction.reply({ content: "Ce sondage est deja termine.", ephemeral: true });
    return;
  }

  poll.ended = true;
  setPoll(key, poll);

  const channel = await interaction.guild.channels.fetch(poll.channelId).catch(() => null);
  if (!channel?.isTextBased()) {
    await interaction.reply({ content: "Salon d'origine introuvable.", ephemeral: true });
    return;
  }

  const message = await channel.messages.fetch(poll.messageId).catch(() => null);
  if (!message) {
    await interaction.reply({ content: "Message du sondage introuvable.", ephemeral: true });
    return;
  }

  const results = await tallyVotes(message, poll.choices);
  const total = results.reduce((a, b) => a + b, 0);
  const embed = buildResultEmbed(poll.question, poll.choices, results, total, interaction.guild, { tag: poll.authorTag }, true);

  await message.edit({ embeds: [embed] });

  addAudit({
    guildId: interaction.guildId,
    action: "poll_end",
    actor: interaction.user.tag,
    target: messageId,
    reason: poll.question,
  });

  await interaction.reply({ content: "Sondage termine. Les resultats ont ete publies.", ephemeral: true });
}

async function handleResults(interaction) {
  const messageId = interaction.options.getString("message_id", true);
  const key = `${interaction.guildId}:${messageId}`;
  const poll = getPoll(key);

  if (!poll) {
    await interaction.reply({ content: "Sondage introuvable. Verifiez l'ID du message.", ephemeral: true });
    return;
  }

  const channel = await interaction.guild.channels.fetch(poll.channelId).catch(() => null);
  if (!channel?.isTextBased()) {
    await interaction.reply({ content: "Salon d'origine introuvable.", ephemeral: true });
    return;
  }

  const message = await channel.messages.fetch(poll.messageId).catch(() => null);
  if (!message) {
    await interaction.reply({ content: "Message du sondage introuvable.", ephemeral: true });
    return;
  }

  const results = await tallyVotes(message, poll.choices);
  const total = results.reduce((a, b) => a + b, 0);
  const embed = buildResultEmbed(poll.question, poll.choices, results, total, interaction.guild, { tag: poll.authorTag }, poll.ended);

  await interaction.reply({ embeds: [embed], ephemeral: true });
}

async function tallyVotes(message, choices) {
  const results = new Array(choices.length).fill(0);
  for (let i = 0; i < choices.length; i++) {
    const reaction = message.reactions.cache.get(EMOJIS[i]);
    if (reaction) {
      const users = await reaction.users.fetch().catch(() => null);
      if (users) {
        results[i] = users.filter((u) => !u.bot).size;
      }
    }
  }
  return results;
}
