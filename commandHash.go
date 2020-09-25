package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func hashFile(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}

	hasher := sha512.New()
	_, err = io.Copy(hasher, f)
	if err != nil {
		return nil
	}

	f.Close()
	return hasher.Sum(nil)
}

func init() {
	h := fmt.Sprintf("%X", hashFile(os.Args[0]))
	if len(h) != 128 {
		panic("Couldn't calculate own hash")
	}
	h = strings.Join(
		[]string{
			"```",
			h[0:32],
			h[32:64],
			h[64:96],
			h[96:128],
			"```",
		},
		"\n",
	)

	mustAddExplicitCommand(&explicitCommand{
		adminOnly:   false,
		chatType:    chatTypeAny,
		command:     "hash",
		helpMessage: "Return hash of the executable",
		function: func(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
			return h, nil
		},
	})
}
