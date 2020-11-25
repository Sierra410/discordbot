package main

<<<<<<< Updated upstream
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
=======
// import (
// 	"fmt"

// 	"discordbot/table"

// 	"github.com/bwmarrin/discordgo"
// )

// func init() {
// 	mustAddExplicitCommand(&explicitCommand{
// 		permLevel:   botPermBotAdmin,
// 		chatType:    chatTypeAny,
// 		command:     "adminctl",
// 		helpMessage: "Usage:\n    add [userid]\n    del [userid]\n    list",
// 		function:    commandAdminctl,
// 	})
// }

// func commandAdminctl(self *explicitCommand, session *discordgo.Session, cmd *parsedCommand) (string, error) {
// 	act := cmd.nthArgOr(0, "")

// 	switch act {
// 	case "add", "del", "delete", "remove":
// 		id := clearUserID(cmd.nthArgOr(1, ""))
// 		if !validateUserId(session, id) {
// 			return "Invalid ID", nil
// 		}

// 		if id == cmd.message.Author.ID {
// 			if act == "add" {
// 				return "No, you cannot _add youself._", nil
// 			} else {
// 				return "You cannot remove yourself", nil
// 			}
// 		}

// 		// Explicit locking and unlocking due to UnsafeSet and UnsafeGet usage
// 		cfg := cfg.Lock(true)
// 		defer cfg.Drop()

// 		admins := interfaceToStringSlice(cfg.Get(cmd.message.GuildID, configAdmins))

// 		idsPermLevel := cfg.GetPermissionLevel(session, cmd.message.GuildID, id)

// 		if act == "add" {
// 			if idsPermLevel > botPermNone {
// 				return "Already an admin", nil
// 			}

// 			admins = append(admins, id)

// 			cfg.Set(cmd.message.GuildID, configAdmins, admins)

// 			return "Added " + idToMention(id), nil
// 		} else {
// 			if idsPermLevel == botPermNone {
// 				return "Not an admin", nil
// 			}

// 			if idsPermLevel > cmd.permLevel {
// 				return "Cannot remove this admin", nil
// 			}

// 			admins, _ := removeFromStringSlice(admins, id)
// 			cfg.Set(cmd.message.GuildID, configAdmins, admins)

// 			return "Removed " + idToMention(id), nil
// 		}
// 	case "list":
// 		admins := interfaceToStringSlice(
// 			cfg.Get(cmd.message.GuildID, configAdmins),
// 		)

// 		table := table.New()
// 		table.Fmtfunc = func(k, v, p string) string {
// 			return fmt.Sprintf("**``%s%s    ``**``%s``", k, p, v)
// 		}

// 		for _, id := range admins {
// 			user, err := session.User(id)
// 			if err != nil {
// 				table.AddLine("ERROR", id)
// 			} else {
// 				table.AddLine(user.Username+"#"+user.Discriminator, id)
// 			}
// 		}

// 		return table.String(), nil
// 	}

// 	return self.helpMessage, nil
// }
>>>>>>> Stashed changes
