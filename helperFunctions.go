package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Discord Stuff

var (
	logChannelName   = "log"
	errInvalidUserId = errors.New("Invalid user id")
)

type channelSearchFunc func(*discordgo.Channel) bool

func sendMultiMessage(session *discordgo.Session, channelId string, s string) {
	pages := pagify(s)

	if s == "" {
		return
	}

	if len(pages) == 1 {
		session.ChannelMessageSend(channelId, pages[0])
		return
	}

	l := strconv.Itoa(len(pages))
	for i, x := range pages {
		session.ChannelMessageSend(channelId, "["+strconv.Itoa(i)+"/"+l+"]\n"+x)
	}
}

func sendToAllChannelsWithName(session *discordgo.Session, guilds interface{}, channelName, message string) {
	sendToGuild := func(guildId string) {
		cs, err := session.GuildChannels(guildId)
		if err != nil {
			return
		}

		c := findChannelByName(cs, channelName)
		if c == nil {
			return
		}

		sendMultiMessage(session, c.ID, message)
	}

	switch guilds.(type) {
	case []*discordgo.Guild:
		for _, guild := range guilds.([]*discordgo.Guild) {
			sendToGuild(guild.ID)
		}
	case []*discordgo.UserGuild:
		for _, guild := range guilds.([]*discordgo.UserGuild) {
			sendToGuild(guild.ID)
		}
	default:
		return
	}
}

// Search helpers

func findChannelByName(channels []*discordgo.Channel, name string) *discordgo.Channel {
	for _, x := range channels {
		if x.Name == name {
			return x
		}
	}
	return nil
}

func findChannel(channels []*discordgo.Channel, sfunc channelSearchFunc) *discordgo.Channel {
	for _, x := range channels {
		if sfunc(x) {
			return x
		}
	}
	return nil
}

func validateUserId(session *discordgo.Session, userid string) bool {
	_, err := session.User(clearUserID(userid))
	return err == nil
}

func validateUserIds(session *discordgo.Session, userids []string) ([]string, error) {
	valid := make([]string, 0, len(userids))

	for _, id := range userids {
		if validateUserId(session, id) {
			valid = append(valid, id)
		}
	}

	return valid, nil
}

// Removes <@ and >, if present
func clearUserID(s string) string {
	if s[0:2] == "<@" && s[len(s)-1] == '>' {
		s = s[2 : len(s)-1]
	}

	return s
}

// Generic stuff

// This is, somehow, faster than append(a, b...) and does less allocations
func joinSlicesOfStrings(a []string, b []string) []string {
	new := make([]string, len(a)+len(b))

	for i := 0; i < len(a); i++ {
		new[i] = a[i]
	}

	for i := 0; i < len(b); i++ {
		new[i+len(a)] = b[i]
	}

	return new
}

func isStringInSlice(s string, sl []string) bool {
	for _, x := range sl {
		if x == s {
			return true
		}
	}

	return false
}

func firstWord(s string) string {
	if s == "" {
		return ""
	}

	start := 0
	for i, c := range s {
		if c != ' ' {
			start = i
			break
		}
	}

	for i, c := range s[start:] {
		if c == ' ' {
			return s[start : start+i]
		}
	}

	return s[start:]
}

func mapDeepcopy(from, to map[string]interface{}) {
	for k, v := range from {
		switch v.(type) {
		case map[string]interface{}:
			m := map[string]interface{}{}
			mapDeepcopy(v.(map[string]interface{}), m)
			to[k] = m
		default:
			to[k] = v
		}
	}
}

// d = (Hour - Second), f = Hour, Minute, Second -> []int64{0, 59, 59}
func durationToParts(d time.Duration, f ...time.Duration) []int64 {
	r := make([]int64, len(f))
	for i, x := range f {
		p := d / x
		d -= p * x
		r[i] = int64(p)
	}

	return r
}

const messagePageSize = 1900

func pagify(s string) []string {
	max := messagePageSize
	offset := 0
	rs := []rune(s)
	pages := []string{}

outer:
	for {
		if offset+max >= len(rs) {
			pages = append(pages, string(rs[offset:]))
			break
		}

		// Try to split at newlines or spaces/tabs
		for _, sep := range []rune{'\n', ' ', '\t'} {
			for i := offset + max; i > offset; i-- {
				if rs[i] == sep {
					pages = append(pages, string(rs[offset:i]))
					offset = i
					continue outer
				}
			}
		}

		// Just split at max
		pages = append(pages, string(rs[offset:offset+max]))
		offset = offset + max
	}

	return pages
}
