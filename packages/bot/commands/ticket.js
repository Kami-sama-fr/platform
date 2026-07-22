import { SlashCommandBuilder } from "discord.js";
import {
  closeSupportTicket,
  getSupportTicket,
} from "../utils/support-tickets.js";
import { sendSupportPanel } from "./support.js";
export const data=new SlashCommandBuilder().setName("ticket").setDescription("Publie le panneau de ticket ou ferme un ticket privé.").addSubcommand(s=>s.setName("create").setDescription("Publier le panneau de création de ticket")).addSubcommand(s=>s.setName("close").setDescription("Fermer le salon actuel"));
export async function execute(i){const sub=i.options.getSubcommand();if(sub==="create"){await sendSupportPanel(i.channel);return i.reply({content:"Le panneau de création de ticket a été publié dans ce salon.",ephemeral:true});}const t=getSupportTicket(i.channelId);if(!t)return i.reply({content:"Ce salon n’est pas un ticket géré par le bot.",ephemeral:true});await i.reply({content:"Fermeture du ticket…",ephemeral:true});closeSupportTicket(i.channelId);setTimeout(()=>i.channel.delete().catch(()=>{}),1000);}
