package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func init() {
	mustAddExplicitCommand(&explicitCommand{
		permLevel:   botPermBotAdmin,
		chatType:    chatTypeAny,
		command:     "adminctl",
		helpMessage: "Usage:\n    add [userid]\n    del [userid]\n    list",
		function:    commandAdminctl,
	})
}

func commandAdminctl(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
	act := cmd.nthArgOr(0, "")

	switch act {
	case "add", "del", "delete", "remove":
		id := clearUserID(cmd.nthArgOr(1, ""))
		if !validateUserId(session, id) {
			return "Invalid ID", nil
		}

		if id == cmd.message.Author.ID {
			if act == "add" {
				return "No, you cannot _add youself._", nil
			} else {
				return "You cannot remove yourself", nil
			}
		}

		// Explicit locking and unlocking due to UnsafeSet and UnsafeGet usage
		cfg := cfg.Lock(true)
		defer cfg.Drop()

		admins := interfaceToStringSlice(cfg.Get(cmd.message.GuildID, configAdmins))

		idsPermLevel := cfg.GetPermissionLevel(session, cmd.message.GuildID, id)

		if act == "add" {
			if idsPermLevel > botPermNone {
				return "Already an admin", nil
			}

			admins = append(admins, id)

			cfg.Set(cmd.message.GuildID, configAdmins, admins)

			return "Added " + idToMention(id), nil
		} else {
			if idsPermLevel == botPermNone {
				return "Not an admin", nil
			}

			if idsPermLevel > cmd.permLevel {
				return "Cannot remove this admin", nil
			}

			admins, _ := removeFromStringSlice(admins, id)
			cfg.Set(cmd.message.GuildID, configAdmins, admins)

			return "Removed " + idToMention(id), nil
		}
	case "list":
		admins := interfaceToStringSlice(
			cfg.Get(cmd.message.GuildID, configAdmins),
		)

		table := NewTable()
		table.Fmtfunc = func(k, v, p string) string {
			return fmt.Sprintf("**``%s%s    ``**``%s``", k, p, v)
		}

		for _, id := range admins {
			user, err := session.User(id)
			if err != nil {
				table.AddLine("ERROR", id)
			} else {
				table.AddLine(user.Username+"#"+user.Discriminator, id)
			}
		}

		return table.String(), nil
	}

	return self.helpMessage, nil
}
