package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func init() {
	mustAddExplicitCommand(&explicitCommand{
		permLevel:   botPermNone,
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
		{"User", "<@" + u.ID + ">"},
		{
			"Cur.",
			"**``" + strings.ReplaceAll(u.Username, "`", "\\`") +
				"#" + u.Discriminator + "``**",
		},
		{"ID", u.ID},
		{"Reg", acDate.Format("**2006-01-02 15:04:05 UTC**")},
		{"Age", acAgeStr},
	}

	resp, err := http.Get(u.AvatarURL("2048"))
	if err == nil {
		defer resp.Body.Close()
	}

	pfpFile := &discordgo.File{
		ContentType: "image/png",
		Name:        u.ID + ".png",
		Reader:      resp.Body,
	}

	dumpMsg, err := session.ChannelMessageSendComplex(
		"758930968849154091",
		&discordgo.MessageSend{
			Files: []*discordgo.File{pfpFile},
		},
	)

	var embedImage *discordgo.MessageEmbedImage
	if err == nil && len(dumpMsg.Attachments) > 0 {
		embedImage = &discordgo.MessageEmbedImage{
			URL: dumpMsg.Attachments[0].URL,
		}
	}

	session.ChannelMessageSendEmbed(
		cmd.message.ChannelID,
		&discordgo.MessageEmbed{
			Description: table.String(),
			Image:       embedImage,
		},
	)

	return "", nil
}

// var pfpFile *discordgo.File
// pfpImage, err := session.UserAvatarDecode(u)

// logDebug.Println("A")

// if err == nil {
// 	r, w := io.Pipe()

// 	go func() {
// 		png.Encode(w, pfpImage)
// 		w.Close()
// 	}()

// 	pfpFile = &discordgo.File{
// 		ContentType: "image/png",
// 		Name:        u.ID + ".png",
// 		Reader:      r,
// 	}
// }

// logDebug.Println("B")

// dumpMsg, err := session.ChannelMessageSendComplex(
// 	"758930968849154091",
// 	&discordgo.MessageSend{
// 		Files: []*discordgo.File{pfpFile},
// 	},
// )

// logDebug.Println("C")

// pfpUrl := ""
// if err == nil {
// 	pfpUrl = dumpMsg.Attachments[0].URL
// } else {
// 	// Just in case
// 	pfpUrl = u.AvatarURL("2048")
// }

// session.ChannelMessageSendEmbed(
// 	cmd.message.ChannelID,
// 	&discordgo.MessageEmbed{
// 		Description: table.String(),
// 		Image: &discordgo.MessageEmbedImage{
// 			URL: pfpUrl,
// 		},
// 	},
// )

// return "", nil
