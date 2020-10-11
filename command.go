package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	chatTypeAny = iota
	chatTypeDm
	chatTypeServer
)

var (
	readyCallbacks   = []func(*discordgo.Session, *discordgo.Ready){}
	explicitCommands = map[string]*explicitCommand{}
	implicitCommands = []*implicitCommand{}

	errNotCommand      = errors.New("Message is not a command")
	errArgsParseFailed = errors.New("Could not parse command arguments")

	errCommandNotFound = errors.New("Command not found")
	errAccessDenied    = errors.New("Access denied")

	errSpaceInCommandName   = errors.New("Command name cannot contain spaces")
	errCommandAlreadyExists = errors.New("Command already exists!")
)

type explicitCommand struct {
	adminOnly           bool
	chatType            int
	command             string
	accessDeniedMessage string // None, if empty
	helpMessage         string
	// If a non-empty string is returned it's sent to the same channel as the command
	// messages > 2K characters are split into several messages
	// if "{help}" is returned explicitCommand.helpMessage is sent instead
	function func(*explicitCommand, *discordgo.Session, *parsedCommand) (string, error)
}

type parsedCommand struct {
	command  string
	args     []string
	body     string
	chatType int
	isAdmin  bool
	message  *discordgo.Message
	prefix   string
}

// Parses a string (message content)
// args are parsed as CSV with " " (space) instead of a comma
func parseCommand(session *discordgo.Session, msg *discordgo.Message) (*parsedCommand, error) {
	prefix := cfg.GetPrefix(msg.GuildID)

	if !strings.HasPrefix(msg.Content, prefix) {
		return nil, errNotCommand
	}

	cmd := &parsedCommand{
		body:    "",
		args:    []string{},
		isAdmin: cfg.IsAdmin(session, msg.GuildID, msg.Author.ID),
		message: msg,
		prefix:  prefix,
	}

	if msg.GuildID == "" {
		cmd.chatType = chatTypeDm
	} else {
		cmd.chatType = chatTypeServer
	}

	lines := strings.SplitN(msg.Content, "\n", 2)
	if len(lines) > 1 {
		cmd.body = lines[1]
	}

	fields := strings.SplitN(lines[0], " ", 2)
	if len(fields) > 1 {
		r := csv.NewReader(strings.NewReader(fields[1]))
		r.Comma = ' '
		r.TrimLeadingSpace = true

		var err error
		cmd.args, err = r.Read()
		if err != nil {
			return nil, errArgsParseFailed
		}
	}

	cmd.command = strings.ToLower(strings.TrimPrefix(fields[0], prefix))

	return cmd, nil
}

// Returns Nth arg or a default value
func (self *parsedCommand) nthArgOr(index int, def string) string {
	if index >= len(self.args) {
		return def
	}
	return self.args[index]
}

func (self *parsedCommand) hasArg(s string) bool {
	for i := range self.args {
		if self.args[i] == s {
			return true
		}
	}
	return false
}

func (self *parsedCommand) execute(session *discordgo.Session) error {
	var err error

	c, ok := explicitCommands[self.command]
	if !ok {
		return errCommandNotFound
	}

	if !self.hasAccess(c) {
		sendMultiMessage(session, self.message.ChannelID, c.accessDeniedMessage)
		return errAccessDenied
	}

	reply, err := c.function(c, session, self)
	if err != nil {
		sendMultiMessage(
			session,
			self.message.ChannelID,
			fmt.Sprintf("Error:\n```%s```", err.Error()),
		)
		return err
	}

	sendMultiMessage(session, self.message.ChannelID, reply)
	return nil
}

func (self *parsedCommand) hasAccess(ec *explicitCommand) bool {
	return !((ec.chatType != 0 && ec.chatType != self.chatType) ||
		(ec.adminOnly && !self.isAdmin))
}

type implicitCommand struct {
	adminOnly bool
	function  func(*discordgo.Session, *discordgo.Message)
}

func addReadyCallback(f func(*discordgo.Session, *discordgo.Ready)) {
	readyCallbacks = append(
		readyCallbacks,
		f,
	)
}

func mustAddExplicitCommand(c *explicitCommand) {
	err := addExplicitCommand(c)
	switch err {
	case errCommandAlreadyExists:
		fmt.Printf(
			"Command \"%s\" already exists!\n",
			strings.ToLower(c.command),
		)

		fallthrough
	default:
		panic(err)
	case nil:
	}
}

func addExplicitCommand(c *explicitCommand) error {
	cmdName := strings.ToLower(c.command)

	if strings.Index(cmdName, " ") != -1 {
		return errSpaceInCommandName
	}

	_, ok := explicitCommands[cmdName]
	if ok {
		return errCommandAlreadyExists
	}

	explicitCommands[cmdName] = c

	return nil
}

func addImplicitCommand(c *implicitCommand) {
	implicitCommands = append(
		implicitCommands,
		c,
	)
}
