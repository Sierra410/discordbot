package main

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	mustAddExplicitCommand(&explicitCommand{
		adminOnly:   false,
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

	if cmd.chatType == chatTypeDm && cmd.isAdmin {
		oldPrefix := cfg.DefaultCommandPrefix
		cfg.DefaultCommandPrefix = newPrefix
		cfg.Save()

		return "Command prefix was changed from ``" + oldPrefix + "`` to ``" + newPrefix + "``", nil
	} else if cmd.chatType == chatTypeServer {

	}

	return "", nil
}
