package main

// import (
// 	"github.com/bwmarrin/discordgo"
// )

// func init() {
// 	mustAddExplicitCommand(&explicitCommand{
// 		permLevel:   botPermServerAdmin,
// 		chatType:    chatTypeAny,
// 		command:     "isadmin",
// 		helpMessage: "Is user an admin?",
// 		function:    commandIsAnAdmin,
// 	})
// }

// func commandIsAnAdmin(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
// 	permLevel := cmd.permLevel

// 	id := clearUserID(cmd.nthArgOr(0, ""))
// 	if id != "" {
// 		permLevel = cfg.GetPermissionLevel(session, cmd.message.GuildID, id)
// 	}

// 	return []string{
// 		"No",
// 		"Yes",
// 		"Yes [Server Owner]",
// 		"Yes [Bot Admin]",
// 		"**Yes [Bot Owner]**",
// 	}[int(permLevel)], nil
// }
