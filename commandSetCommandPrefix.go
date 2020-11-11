package main

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	mustAddExplicitCommand(&explicitCommand{
		permLevel:   botPermServerAdmin,
		chatType:    chatTypeAny,
		command:     "setprefix",
		helpMessage: "Usage: setprefix newprefix",
		function:    commandSetCommandPrefix,
	})
}

var (
	errPrefixCannotContainSpaces = errors.New("Command prefix cannot contain spaces!")
)

func commandSetCommandPrefix(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
	newPrefix := cmd.nthArgOr(0, "")

	if newPrefix == "" {
		return self.helpMessage, nil
	}

	if strings.IndexAny(newPrefix, " \t\n") != -1 {
		return errPrefixCannotContainSpaces.Error(), nil
	}

	cfg.Set(cmd.message.GuildID, configCommandPrefix, newPrefix)

	return "Command prefix was changed from ``" + cmd.prefix + "`` to ``" + newPrefix + "``", nil
}
