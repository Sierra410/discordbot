package bot

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
		isa := cfg.IsAdmin(session, msg.GuildID, msg.Author.ID)
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
