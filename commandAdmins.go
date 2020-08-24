package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	mustAddExplicitCommand(&explicitCommand{
		adminOnly:   true,
		chatType:    chatTypeDm,
		command:     "adminctl",
		helpMessage: "Usage:\n    add [userid...]\n    del [userid...]\n    list",
		function:    commandAdminctl,
	})
}

func commandAdminctl(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
	act := cmd.nthArgOr(0, "")

	switch act {
	case "add", "del":
		ids, err := validateUserIds(session, cmd.args[1:])
		if err != nil {
			return "", errInvalidUserId
		}

		if act == "add" {
			// add
			newAdmins := cfg.Admins

			for _, x := range ids {
				if !isStringInSlice(x, newAdmins) {
					newAdmins = append(newAdmins, x)
				}
			}

			cfg.Admins = newAdmins
		} else {
			// del
			newAdmins := []string{}

			for _, x := range cfg.Admins {
				if !isStringInSlice(x, ids) {
					newAdmins = append(newAdmins, x)
				}
			}

			cfg.Admins = newAdmins
		}

		err = cfg.Save()
		if err != nil {
			return "", err
		}
	case "list":
		// To do: Username. Long list safety.
		b := strings.Builder{}
		for _, x := range cfg.Admins {
			b.WriteString(x)
			b.WriteRune('\n')
		}

		return b.String(), nil
	}

	return self.helpMessage, nil
}
