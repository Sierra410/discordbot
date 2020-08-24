package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	mustAddExplicitCommand(&explicitCommand{
		adminOnly:   false,
		chatType:    chatTypeAny,
		command:     "userprobe",
		helpMessage: "Usage: userprobe userid",
		function:    commandUserprobe,
	})
}

func commandUserprobe(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
	id := cmd.nthArgOr(0, cmd.message.Author.ID)

	u, err := session.User(clearUserID(id))
	if err != nil {
		return "", err
	}

	acDate, err := discordgo.SnowflakeTimestamp(u.ID)
	if err != nil {
		return "", err
	}
	acDate = acDate.UTC()
	acAge := durationToParts(
		time.Now().Sub(acDate).Round(time.Second),
		(time.Hour * 24),
		time.Hour,
		time.Minute,
		time.Second,
	)

	acAgeStr := fmt.Sprintf(
		"**%d Days %02d:%02d:%02d**",
		acAge[0],
		acAge[1],
		acAge[2],
		acAge[3],
	)

	table := NewTable()
	table.Fmtfunc = func(k, v, p string) string {
		return "**`` " + k + p + " ``**  " + v
	}

	table.Lines = [][2]string{
		[2]string{"User", "<@" + u.ID + ">"},
		[2]string{
			"Cur.",
			"**``" + strings.ReplaceAll(u.Username, "`", "\\`") +
				"#" + u.Discriminator + "``**",
		},
		[2]string{"ID", u.ID},
		[2]string{"Reg", acDate.Format("**2006-01-02 15:04:05 UTC**")},
		[2]string{"Age", acAgeStr},
	}

	session.ChannelMessageSendEmbed(
		cmd.message.ChannelID,
		&discordgo.MessageEmbed{
			Description: table.String(),
			Image: &discordgo.MessageEmbedImage{
				URL: u.AvatarURL("2048"),
			},
		},
	)

	return "", nil
}
