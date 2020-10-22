package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func init() {
	addReadyCallback(matrixChatWatchdogLoop)
}

func matrixChatWatchdogLoop(session *discordgo.Session, data *discordgo.Ready) {
	prefix := "matrix"
	homeserver := cfg.GetOr(prefix, "homeserver", "matrix.org").(string)
	username := cfg.Get(prefix, "username").(string)
	password := cfg.Get(prefix, "password").(string)
	channels := cfg.Get(prefix, "channels").(map[string]interface{})

	for {
		logErr.Println(
			matrixChatWatchdog(session, homeserver, username, password, channels),
			"\nRestarting Matrix client in 10 seconds...",
		)
		time.Sleep(time.Second * 10)
	}
}

func matrixChatWatchdog(session *discordgo.Session, homeserver, username, password string, channels map[string]interface{}) error {
	fmt.Printf("MATRIX: %s@%s\n", username, homeserver)

	client, err := mautrix.NewClient(homeserver, "", "")
	if err != nil {
		return err
	}

	_, err = client.Login(
		&mautrix.ReqLogin{
			Type: "m.login.password",
			Identifier: mautrix.UserIdentifier{
				Type: mautrix.IdentifierTypeUser,
				User: username,
			},
			Password:         password,
			StoreCredentials: true,
		},
	)
	if err != nil {
		return err
	}

	syncer := client.Syncer.(*mautrix.DefaultSyncer)

	var f = func(source mautrix.EventSource, evt *event.Event) {
		age := eventAge(evt.Timestamp)
		reportChannel, _ := channels[string(evt.RoomID)].(string)

		if age < 30 && reportChannel != "" {
			session.ChannelMessageSend(
				reportChannel,
				"New message in Matrix from "+string(evt.Sender),
			)
		}
	}

	syncer.OnEventType(
		event.EventMessage,
		f,
	)

	syncer.OnEventType(
		event.EventEncrypted,
		f,
	)

	return client.Sync()
}

func eventAge(timestamp int64) int64 {
	return time.Now().Unix() - (timestamp / 1000)
}
