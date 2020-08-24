package main

import "github.com/bwmarrin/discordgo"

func init() {
	addReadyCallback(func(s *discordgo.Session, r *discordgo.Ready) {
		updateStatus(s)
	})

	mustAddExplicitCommand(&explicitCommand{
		adminOnly:   true,
		chatType:    chatTypeDm,
		command:     "setstatus",
		helpMessage: "Usage: setstatus g/l/s/w\n    status",
		function:    commandStatus,
	})
}

func commandStatus(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
	newStatus := cmd.body
	gameType := "g"

	for _, x := range cmd.args {
		if len(x) != 1 {
			continue
		}

		gt := gameTypeAtoi(rune(x[0]))
		if gt != -1 {
			gameType = x
			break
		}
	}

	cfg.Status = gameType + newStatus
	cfg.Save()

	updateStatus(session)

	return "Ok", nil
}

func updateStatus(session *discordgo.Session) {
	status := ""
	gameType := discordgo.GameTypeGame

	if len(cfg.Status) >= 2 {
		status = cfg.Status[1:]
		gameType = gameTypeAtoi(rune(cfg.Status[0]))
		if gameType == -1 {
			gameType = discordgo.GameTypeGame
		}
	}

	session.UpdateStatusComplex(
		discordgo.UpdateStatusData{
			Game: &discordgo.Game{
				Name: status,
				Type: gameType,
			},
		},
	)
}

func gameTypeAtoi(r rune) discordgo.GameType {
	switch r {
	case 'g':
		return discordgo.GameTypeGame
	case 'l':
		return discordgo.GameTypeListening
	case 's':
		return discordgo.GameTypeStreaming
	case 'w':
		return discordgo.GameTypeWatching
	default:
		return -1
	}
}
