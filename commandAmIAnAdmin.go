package main

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	mustAddExplicitCommand(&explicitCommand{
		adminOnly:   false,
		chatType:    chatTypeAny,
		command:     "amianadmin",
		helpMessage: "Tells you that",
		function:    commandAmIAnAdmin,
	})
}

func commandAmIAnAdmin(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
	admins := interfaceToInterfaceSlice(cfg.Get("", configAdmins))
	if isInSlice(cmd.message.Author.ID, admins) {
		return "**Yes [Bot Owner]**", nil
	}

	admins = interfaceToInterfaceSlice(cfg.Get(cmd.message.GuildID, configAdmins))
	if isInSlice(cmd.message.Author.ID, admins) {
		return "Yes [Bot Admin]", nil
	}

	guild, err := session.Guild(cmd.message.GuildID)
	if err == nil && guild.OwnerID == cmd.message.Author.ID {
		return "Yes [Server Owner]", nil
	}

	return "No", nil
}
