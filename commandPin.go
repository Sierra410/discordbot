package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	mustAddExplicitCommand(&explicitCommand{
		permLevel:   botPermServerAdmin,
		chatType:    chatTypeAny,
		command:     "pin",
		helpMessage: "Usage: list to list messages.\nReply to a message with this command to pin/unpin it.",
		function:    commandPinMessage,
	})
}

func getMessageLinkFromRef(mr *discordgo.MessageReference) string {
	if mr == nil {
		return ""
	}

	return getMessageLink(mr.GuildID, mr.ChannelID, mr.MessageID)
}

func getMessageLink(guildId, channelId, messageId string) string {
	if guildId == "" || channelId == "" || messageId == "" {
		return ""
	}

	return fmt.Sprintf("https://discord.com/channels/%s/%s/%s", guildId, channelId, messageId)
}

func commandPinMessage(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
	guildId := cmd.message.GuildID
	channelId := cmd.message.ChannelID
	messageId := cmd.message.ID

	cfg := cfg.Lock(true)
	defer cfg.Drop()

	cfgKey := channelId + "_pinned"

	pinned := interfaceToStringSlice(
		cfg.GetOr(guildId, cfgKey, []string{}),
	)

	if cmd.nthArgOr(0, "") == "list" {
		b := strings.Builder{}
		b.WriteString("Pinned messages:\n")
		for _, x := range pinned {
			b.WriteString("<")
			b.WriteString(getMessageLink(guildId, channelId, x))
			b.WriteString(">\n")
		}

		return b.String(), nil
	}

	mr := cmd.message.MessageReference
	if mr == nil {
		return self.helpMessage, nil
	}

	if mr.ChannelID != channelId || mr.GuildID != guildId {
		return "I can't pin _that_ O\\_o", nil
	}

	if !isStringInSlice(mr.MessageID, pinned) {
		//add
		pinned = append(pinned, mr.MessageID)

		cfg.Set(guildId, cfgKey, pinned)
		return "Pinned " + messageId, nil
	} else {
		//remove
		pinned, removed := removeFromStringSlice(pinned, mr.MessageID)

		if len(removed) == 0 {
			return "Are you sure this is pinned?", nil
		}

		cfg.Set(guildId, cfgKey, pinned)
		return "Removed " + messageId, nil
	}
}
