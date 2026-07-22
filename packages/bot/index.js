import { readdir } from "node:fs/promises";
import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import { Client, Collection, GatewayIntentBits, Partials } from "discord.js";
import { env } from "./config/env.js";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMembers,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.MessageContent,
    GatewayIntentBits.GuildModeration,
    GatewayIntentBits.GuildMessageReactions,
    GatewayIntentBits.GuildVoiceStates,
    GatewayIntentBits.GuildScheduledEvents,
    GatewayIntentBits.AutoModerationConfiguration,
    GatewayIntentBits.AutoModerationExecution,
    GatewayIntentBits.DirectMessages,
  ],
  partials: [
    Partials.Channel,
    Partials.Message,
    Partials.Reaction,
    Partials.User,
    Partials.GuildMember,
  ],
});

client.commands = new Collection();

async function loadCommands() {
  const commandsPath = path.join(__dirname, "commands");
  const commandFiles = (await readdir(commandsPath)).filter((file) => file.endsWith(".js"));
  for (const file of commandFiles) {
    const command = await import(pathToFileURL(path.join(commandsPath, file)).href);
    if (!command.data || !command.execute) {
      console.warn(`[WARNING] La commande ${file} n'exporte pas data + execute.`);
      continue;
    }
    client.commands.set(command.data.name, command);
  }
}

async function loadEvents() {
  const eventsPath = path.join(__dirname, "events");
  const eventFiles = (await readdir(eventsPath)).filter((file) => file.endsWith(".js"));
  for (const file of eventFiles) {
    const event = await import(pathToFileURL(path.join(eventsPath, file)).href);
    if (event.once) client.once(event.name, (...args) => event.execute(...args));
    else client.on(event.name, (...args) => event.execute(...args));
  }
}

async function bootstrap() {
  try {
    if (!env.isConfigured) {
      console.warn(`[WARN] Discord bot config incomplete. Missing: ${env.missingRequiredVars.join(", ")}`);
      setInterval(() => {}, 1 << 30);
      return;
    }
    await loadCommands();
    await loadEvents();
    await client.login(env.token);
  } catch (error) {
    console.error("Échec du démarrage du bot:", error);
    process.exit(1);
  }
}
void bootstrap();
