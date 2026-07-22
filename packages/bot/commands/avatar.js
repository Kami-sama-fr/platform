import { SlashCommandBuilder, EmbedBuilder } from "discord.js";
export const data = new SlashCommandBuilder().setName("avatar").setDescription("Affiche l’avatar d’un utilisateur.").addUserOption(o=>o.setName("user").setDescription("Utilisateur ciblé"));
export async function execute(i){const u=i.options.getUser("user")??i.user;const url=u.displayAvatarURL({size:4096});await i.reply({embeds:[new EmbedBuilder().setTitle(`Avatar de ${u.username}`).setImage(url).setURL(url)]});}
