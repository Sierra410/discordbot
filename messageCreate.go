package main

import (
	"github.com/bwmarrin/discordgo"
)

func messageCreate(session *discordgo.Session, data *discordgo.MessageCreate) {
	msg := data.Message

	if msg.Author.ID == ownId {
		return
	}

	prefix, ok := cfg.GetPerServerConfig(data.GuildID)["CommandPrefix"].(string)
	if !ok {
		prefix = cfg.DefaultCommandPrefix
	}

	cmd, err := parseCommand(msg, prefix)
	switch err {
	case nil:
		_ = cmd.execute(session)
	case errArgsParseFailed:
		sendMultiMessage(session, msg.ChannelID, err.Error())
	case errNotCommand:
		isa := cfg.IsAdmin(msg.Author.ID)
		for _, ic := range implicitCommands {
			if ic.adminOnly && !isa {
				continue
			}

			ic.function(session, msg)
		}
	default:
		// Unknown error
	}
}
