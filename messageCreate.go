package main

import (
	"github.com/bwmarrin/discordgo"
)

func messageCreate(session *discordgo.Session, data *discordgo.MessageCreate) {
	msg := data.Message

	if msg.Author.ID == ownId {
		return
	}

	cmd, err := parseCommand(session, msg)
	switch err {
	case nil:
		_ = cmd.execute(session)
	case errArgsParseFailed:
		sendMultiMessage(session, msg.ChannelID, err.Error())
	case errNotCommand:
		for _, ic := range implicitCommands {
			ic.execute(session, msg)
		}
	default:
		// Unknown error
	}
}
