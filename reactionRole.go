package main

// type ReactionRole struct {
// 	ChannelId string
// 	MessageId string
// 	RoleId    string
// 	Emoji     string
// }

// func setMembership(session *discordgo.Session, guildId, userId, roleId string, membership bool) {
// 	if membership {
// 		// Add
// 		session.GuildMemberRoleAdd(
// 			guildId,
// 			userId,
// 			roleId,
// 		)
// 	} else {
// 		// Remove
// 		session.GuildMemberRoleRemove(
// 			guildId,
// 			userId,
// 			roleId,
// 		)
// 	}
// }

// func processReaction(session *discordgo.Session, guildId,
// 	messageId, userId, emojiName string, added bool) bool {

// 	if userId == ownId {
// 		return false
// 	}

// 	rc, ok := cfg.ReactionRoles[guildId]
// 	if !ok {
// 		return false
// 	}

// 	if messageId != rc.MessageId {
// 		return false
// 	}

// 	if emojiName != rc.Emoji {
// 		return false
// 	}

// 	setMembership(session, guildId, userId, rc.RoleId, added)

// 	return true
// }

// func rulesReactionAdd(session *discordgo.Session, data *discordgo.MessageReactionAdd) {
// 	processReaction(
// 		session,
// 		data.GuildID,
// 		data.MessageID,
// 		data.UserID,
// 		data.Emoji.Name,
// 		true,
// 	)
// }

// func rulesReationRemove(session *discordgo.Session, data *discordgo.MessageReactionRemove) {
// 	processReaction(
// 		session,
// 		data.GuildID,
// 		data.MessageID,
// 		data.UserID,
// 		data.Emoji.Name,
// 		false,
// 	)
// }
