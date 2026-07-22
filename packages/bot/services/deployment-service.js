import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { EmbedBuilder } from "discord.js";
import { env } from "../config/env.js";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const stateDirectory = path.resolve(__dirname, "../data");
const stateFile = path.join(stateDirectory, "deployment-state.json");

export function getReleaseInfo() {
  return {
    version: env.version,
    commitSha: env.commitSha,
    shortCommit: env.commitSha === "unknown" ? "unknown" : env.commitSha.slice(0, 7),
    buildDate: env.buildDate,
    environment: env.environment,
  };
}

async function readState() {
  try {
    return JSON.parse(await fs.readFile(stateFile, "utf8"));
  } catch {
    return {};
  }
}

async function writeState(state) {
  await fs.mkdir(stateDirectory, { recursive: true });
  await fs.writeFile(stateFile, `${JSON.stringify(state, null, 2)}\n`, "utf8");
}

export async function announceDeployment(client) {
  if (!env.announceDeployments || !env.updatesChannelId) return false;

  const release = getReleaseInfo();
  const deploymentKey = `${release.environment}:${release.version}:${release.commitSha}`;
  const state = await readState();

  if (state.lastAnnouncedDeployment === deploymentKey) return false;

  const channel = await client.channels.fetch(env.updatesChannelId).catch(() => null);
  if (!channel?.isTextBased()) {
    console.warn("[DEPLOYMENT] DISCORD_UPDATES_CHANNEL_ID ne référence pas un salon textuel accessible.");
    return false;
  }

  const embed = new EmbedBuilder()
    .setTitle("🚀 Discord Enterprise mis à jour")
    .setDescription("Une nouvelle version du bot vient d’être déployée et ses commandes slash ont été synchronisées.")
    .addFields(
      { name: "Version", value: release.version, inline: true },
      { name: "Commit", value: release.shortCommit, inline: true },
      { name: "Environnement", value: release.environment, inline: true },
      { name: "Build", value: release.buildDate },
    )
    .setTimestamp();

  await channel.send({ embeds: [embed] });
  await writeState({
    lastAnnouncedDeployment: deploymentKey,
    announcedAt: new Date().toISOString(),
  });

  return true;
}
