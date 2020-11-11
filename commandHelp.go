package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	mustAddExplicitCommand(&explicitCommand{
		permLevel:   botPermNone,
		chatType:    chatTypeAny,
		command:     "help",
		helpMessage: "Yes",
		function:    commandHelp,
	})
}

func commandHelp(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
	arg := cmd.nthArgOr(0, "")

	b := strings.Builder{}

	switch arg {
	case "":
		b.WriteString("Available commands:\n")
		for _, v := range explicitCommands {
			if cmd.hasAccess(v) {
				b.WriteString("    ")
				b.WriteString(v.command)
				b.WriteRune('\n')
			}
		}
		return b.String(), nil

	default:
		c, ok := explicitCommands[arg]
		if !ok {
			return errCommandNotFound.Error(), nil
		}

		if c.helpMessage == "" || !cmd.hasAccess(c) {
			return "Help message is not available", nil
		}

		return c.helpMessage, nil
	}
}
