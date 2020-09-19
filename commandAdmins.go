package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	mustAddExplicitCommand(&explicitCommand{
		adminOnly:   true,
		chatType:    chatTypeAny,
		command:     "adminctl",
		helpMessage: "Usage:\n    add [userid...]\n    del [userid...]\n    list",
		function:    commandAdminctl,
	})
}

func commandAdminctl(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
	act := cmd.nthArgOr(0, "")

	switch act {
	case "add", "del", "delete", "remove":
		ids := validateUserIds(session, clearUserIDs(cmd.args[1:]))
		removeFromStringSlice(ids, cmd.message.Author.ID)

		if len(ids) == 0 {
			return "", nil
		}

		// Explicit locking and unlocking due to UnsafeSet and UnsafeGet usage
		cfg.cfg.Lock()
		defer cfg.cfg.Unlock()

		admins := []string{}
		adminsInterface := cfg.UnsafeGet(cmd.message.GuildID, configAdmins)
		switch x := adminsInterface.(type) {
		case []interface{}:
			admins = interfaceSliceToStringSlice(x)
		case []string:
			admins = x
		}

		if act == "add" {
			admins = append(admins, ids...)

			cfg.UnsafeSet(cmd.message.GuildID, configAdmins, admins)

			return "Added:\n" + strings.Join(idsToMentions(ids), "\n"), nil
		} else {
			removed := []string{}
			admins, removed = removeFromStringSlice(admins, ids...)

			cfg.UnsafeSet(cmd.message.GuildID, configAdmins, admins)

			if len(removed) != 0 {
				return "Removed:\n" + strings.Join(idsToMentions(removed), "\n"), nil
			}

			return "", nil
		}
	case "list":
		admins := interfaceSliceToStringSlice(
			interfaceToInterfaceSlice(
				cfg.Get(cmd.message.GuildID, configAdmins),
			),
		)

		b := strings.Builder{}

		b.WriteString("Admins:\n")
		for _, x := range admins {
			b.WriteString(interfaceToString(x))
			b.WriteRune('\n')
		}

		return b.String(), nil
	}

	return self.helpMessage, nil
}
