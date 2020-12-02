package main

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
)

func init() {
	s4re := regexp.MustCompile(`(^|[\s*_~|` + "`" + `])[Ss]4\d{0,2}('s)?([\s*_~|` + "`" + `]|$)`)
	s4id := "648266819618209793"

	addImplicitCommand(&implicitCommand{
		permLevel: botPermNone,
		function: func(session *discordgo.Session, message *discordgo.Message) string {
			if message.Author.Bot || message.Author.ID == s4id {
				return ""
			}

			if s4re.FindStringIndex(message.Content) != nil {
				c, err := session.UserChannelCreate(s4id)
				if err != nil {
					logErr.Println(err)
					return ""
				}

				messageLink := getMessageLinkFromMessage(message)

				logInfo.Println(messageLink)

				session.ChannelMessageSendEmbed(
					c.ID,
					&discordgo.MessageEmbed{
						Description: messageLink + "\n\n" + message.Content,
						Author: &discordgo.MessageEmbedAuthor{
							Name:    message.Author.String(),
							IconURL: message.Author.AvatarURL("128"),
						},
						Footer: &discordgo.MessageEmbedFooter{
							Text: "UserID: " + message.Author.ID,
						},
					},
				)
			}

			return ""
		},
	})
}
