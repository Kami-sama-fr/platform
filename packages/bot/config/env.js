import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import dotenv from "dotenv";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const botRoot = path.resolve(__dirname, "..");
const workspaceRoot = path.resolve(botRoot, "..");
const envFile = path.join(workspaceRoot, ".env");

if (fs.existsSync(envFile)) {
  dotenv.config({ path: envFile, override: false });
}

const requiredVars = ["DISCORD_TOKEN", "DISCORD_CLIENT_ID", "DISCORD_GUILD_ID"];
const missingRequiredVars = requiredVars.filter((key) => !process.env[key]);

export const env = {
  token: process.env.DISCORD_TOKEN ?? "",
  clientId: process.env.DISCORD_CLIENT_ID ?? "",
  guildId: process.env.DISCORD_GUILD_ID ?? "",
  modoRoleId: process.env.DISCORD_MODO_ROLE_ID ?? "",
  supportCategoryId: process.env.DISCORD_SUPPORT_CATEGORY_ID ?? "",
  verifiedRoleId: process.env.DISCORD_ROLE_VERIFIED_ID ?? "",
  memberRoleId: process.env.DISCORD_ROLE_MEMBER_ID ?? "",
  agentsRoleId: process.env.DISCORD_ROLE_AGENTS_ID ?? "",
  updatesChannelId: process.env.DISCORD_UPDATES_CHANNEL_ID ?? "",
  panelChannelId: process.env.DISCORD_PANEL_CHANNEL_ID ?? "",
  levelChannelId: process.env.DISCORD_LEVEL_CHANNEL_ID ?? "",
  rulesChannelId: process.env.DISCORD_RULES_CHANNEL_ID ?? "",
  welcomeChannelId: process.env.DISCORD_WELCOME_CHANNEL_ID ?? "",
  welcomeDmEnabled: process.env.DISCORD_WELCOME_DM_ENABLED !== "false",
  commandScope: process.env.DISCORD_COMMAND_SCOPE ?? "global",
  version: process.env.APP_VERSION ?? "1.1.0",
  commitSha: process.env.GIT_COMMIT_SHA ?? process.env.COMMIT_SHA ?? "unknown",
  buildDate: process.env.BUILD_DATE ?? "unknown",
  environment: process.env.NODE_ENV ?? "development",
  announceDeployments: process.env.DISCORD_ANNOUNCE_DEPLOYMENTS !== "false",
  missingRequiredVars,
  isConfigured: missingRequiredVars.length === 0,
};

export function assertDiscordEnv() {
  if (!env.isConfigured) {
    throw new Error(
      `Missing required environment variables: ${env.missingRequiredVars.join(", ")}`
    );
  }
}
