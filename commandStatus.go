package main

// import "github.com/bwmarrin/discordgo"

// func init() {
// 	addReadyCallback(func(s *discordgo.Session, r *discordgo.Ready) {
// 		updateStatus(s)
// 	})

<<<<<<< Updated upstream
	mustAddExplicitCommand(&explicitCommand{
		adminOnly:   true,
		chatType:    chatTypeDm,
		command:     "setstatus",
		helpMessage: "Usage: setstatus g/l/s/w\n    status",
		function:    commandStatus,
	})
}
=======
// 	mustAddExplicitCommand(&explicitCommand{
// 		permLevel:   botPermBotAdmin,
// 		chatType:    chatTypeAny,
// 		command:     "setstatus",
// 		helpMessage: "Usage: setstatus g/l/s/w\n    status",
// 		function:    commandStatus,
// 	})
// }
>>>>>>> Stashed changes

// func commandStatus(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
// 	newStatus := cmd.body
// 	gameType := "g"

// 	for _, x := range cmd.args {
// 		if len(x) != 1 {
// 			continue
// 		}

// 		gt := gameTypeAtoi(rune(x[0]))
// 		if gt != -1 {
// 			gameType = x
// 			break
// 		}
// 	}

// 	cfg.Set("", configStatus, gameType+newStatus)

// 	updateStatus(session)

// 	return "Ok", nil
// }

// func updateStatus(session *discordgo.Session) {
// 	status := ""
// 	gameType := discordgo.GameTypeGame

// 	tStatus := cfg.Get("", configStatus).(string)
// 	if len(tStatus) >= 2 {
// 		status = tStatus[1:]
// 		gameType = gameTypeAtoi(rune(tStatus[0]))
// 		if gameType == -1 {
// 			gameType = discordgo.GameTypeGame
// 		}
// 	}

// 	session.UpdateStatusComplex(
// 		discordgo.UpdateStatusData{
// 			Game: &discordgo.Game{
// 				Name: status,
// 				Type: gameType,
// 			},
// 		},
// 	)
// }

// func gameTypeAtoi(r rune) discordgo.GameType {
// 	switch r {
// 	case 'g':
// 		return discordgo.GameTypeGame
// 	case 'l':
// 		return discordgo.GameTypeListening
// 	case 's':
// 		return discordgo.GameTypeStreaming
// 	case 'w':
// 		return discordgo.GameTypeWatching
// 	default:
// 		return -1
// 	}
// }
