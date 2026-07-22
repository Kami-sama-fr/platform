import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const dataDir = path.join(__dirname, "..", "data");
const storePath = path.join(dataDir, "store.json");

function ensureStoreFile() {
  if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir, { recursive: true });
  }

  if (!fs.existsSync(storePath)) {
    fs.writeFileSync(
      storePath,
      JSON.stringify(
        {
          warnings: {},
          incidents: {},
          tickets: {},
          auditLog: [],
          xp: {},
          memberEvents: {},
          welcomeSettings: {},
          goodbyeSettings: {},
          panelSettings: {},
          notificationChannels: {},
          polls: {},
        },
        null,
        2
      ),
      "utf8"
    );
  }
}

function loadStore() {
  ensureStoreFile();

  try {
    const raw = fs.readFileSync(storePath, "utf8");
    return JSON.parse(raw);
  } catch (error) {
    console.error("Impossible de charger le store du bot, reinitialisation:", error);
    return {
      warnings: {},
      incidents: {},
      tickets: {},
      auditLog: [],
      xp: {},
      memberEvents: {},
      welcomeSettings: {},
      goodbyeSettings: {},
      panelSettings: {},
      notificationChannels: {},
      polls: {},
    };
  }
}

function persistStore() {
  ensureStoreFile();

  const payload = {
    warnings: Object.fromEntries(warnings),
    incidents: Object.fromEntries(incidents),
    tickets: Object.fromEntries(tickets),
    auditLog,
    xp: Object.fromEntries(xp),
    memberEvents: Object.fromEntries(memberEvents),
    welcomeSettings: Object.fromEntries(welcomeSettings),
    goodbyeSettings: Object.fromEntries(goodbyeSettings),
    panelSettings: Object.fromEntries(panelSettings),
    notificationChannels: Object.fromEntries(notificationChannels),
    polls: Object.fromEntries(polls),
  };

  const tmpPath = `${storePath}.tmp`;
  fs.writeFileSync(tmpPath, JSON.stringify(payload, null, 2), "utf8");
  fs.renameSync(tmpPath, storePath);
}

const initialStore = loadStore();

function reviveDateEntries(entries) {
  return entries.map((entry) => ({
    ...entry,
    at: entry.at ? new Date(entry.at) : new Date(),
  }));
}

export const warnings = new Map(
  Object.entries(initialStore.warnings ?? {}).map(([key, value]) => [key, reviveDateEntries(value)])
);
export const incidents = new Map(Object.entries(initialStore.incidents ?? {}));
export const tickets = new Map(Object.entries(initialStore.tickets ?? {}));
export const auditLog = reviveDateEntries(initialStore.auditLog ?? []);
export const xp = new Map(Object.entries(initialStore.xp ?? {}));
export const memberEvents = new Map(Object.entries(initialStore.memberEvents ?? {}));
export const welcomeSettings = new Map(Object.entries(initialStore.welcomeSettings ?? {}));
export const goodbyeSettings = new Map(Object.entries(initialStore.goodbyeSettings ?? {}));
export const panelSettings = new Map(Object.entries(initialStore.panelSettings ?? {}));
export const notificationChannels = new Map(Object.entries(initialStore.notificationChannels ?? {}));
export const polls = new Map(Object.entries(initialStore.polls ?? {}));

export function areMemberEventsEnabled(guildId) {
  return memberEvents.get(guildId) !== false;
}

export function setMemberEventsEnabled(guildId, enabled) {
  memberEvents.set(guildId, enabled);
  persistStore();
}

export function getWelcomeSettings(guildId) {
  return welcomeSettings.get(guildId) ?? {};
}

export function updateWelcomeSettings(guildId, settings) {
  welcomeSettings.set(guildId, {
    ...getWelcomeSettings(guildId),
    ...settings,
  });
  persistStore();
}

export function resetWelcomeSettings(guildId) {
  welcomeSettings.delete(guildId);
  persistStore();
}

export function getGoodbyeSettings(guildId) {
  return goodbyeSettings.get(guildId) ?? {};
}

export function updateGoodbyeSettings(guildId, settings) {
  goodbyeSettings.set(guildId, {
    ...getGoodbyeSettings(guildId),
    ...settings,
  });
  persistStore();
}

export function resetGoodbyeSettings(guildId) {
  goodbyeSettings.delete(guildId);
  persistStore();
}

export function getPanelSettings(guildId) {
  return panelSettings.get(guildId) ?? {};
}

export function updatePanelSettings(guildId, settings) {
  panelSettings.set(guildId, {
    ...getPanelSettings(guildId),
    ...settings,
  });
  persistStore();
}

export function resetPanelSettings(guildId) {
  panelSettings.delete(guildId);
  persistStore();
}

export function saveStore() {
  persistStore();
}

export function addAudit(entry) {
  auditLog.unshift({ id: crypto.randomUUID(), at: new Date(), ...entry });
  if (auditLog.length > 500) auditLog.length = 500;
  persistStore();
}

const NOTIFICATION_TYPES = ["server_info", "infra_info", "social_network", "status"];

export function getNotificationTypes() {
  return NOTIFICATION_TYPES;
}

export function getNotificationChannels(guildId) {
  return notificationChannels.get(guildId) ?? {};
}

export function setNotificationChannel(guildId, type, channelId) {
  const channels = getNotificationChannels(guildId);
  channels[type] = channelId;
  notificationChannels.set(guildId, channels);
  persistStore();
}

export function removeNotificationChannel(guildId, type) {
  const channels = getNotificationChannels(guildId);
  delete channels[type];
  notificationChannels.set(guildId, channels);
  persistStore();
}

export function resetNotificationChannels(guildId) {
  notificationChannels.delete(guildId);
  persistStore();
}

export function getPoll(key) {
  return polls.get(key) ?? null;
}

export function setPoll(key, data) {
  polls.set(key, data);
  persistStore();
}

export function removePoll(key) {
  polls.delete(key);
  persistStore();
}
