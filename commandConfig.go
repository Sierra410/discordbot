package main

// import "github.com/bwmarrin/discordgo"

<<<<<<< Updated upstream
func init() {
	mustAddExplicitCommand(&explicitCommand{
		adminOnly:   true,
		chatType:    chatTypeDm,
		command:     "config",
		helpMessage: "Usage: config reload",
		function:    commandConfig,
	})
}
=======
// func init() {
// 	mustAddExplicitCommand(&explicitCommand{
// 		permLevel:   botPermBotAdmin,
// 		chatType:    chatTypeAny,
// 		command:     "config",
// 		helpMessage: "Usage: config reload",
// 		function:    commandConfig,
// 	})
// }
>>>>>>> Stashed changes

// func commandConfig(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
// 	switch {
// 	case cmd.hasArg("reload"):
// 		logInfo.Println("Reloading config")
// 		cfg.Load()

// 		return "Reloading config", nil
// 	}

// 	return self.helpMessage, nil
// }
