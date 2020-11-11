package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func init() {
	addReadyCallback(matrixChatWatchdog)
}

func matrixChatWatchdog(session *discordgo.Session, ready *discordgo.Ready) {
	prefix := "matrix"
	homeserver := cfg.GetOr(prefix, "homeserver", "matrix.org").(string)
	username := cfg.Get(prefix, "username").(string)
	password := cfg.Get(prefix, "password").(string)
	channels := map[string]string{}

	for k, v := range cfg.Get(prefix, "channels").(map[string]interface{}) {
		switch v := v.(type) {
		case string:
			channels[k] = v
		}
	}

	fmt.Printf("MATRIX: %s@%s\n", username, homeserver)

	client, err := mautrix.NewClient(homeserver, "", "")
	if err != nil {
		logErr.Println(err)
		return
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
		logErr.Panicln(err)
		return
	}

	// Not scalable! Will work weirdly with more than 1 channel!
	var (
		lock        = sync.Mutex{}
		lastSender  = ""
		sent        = 0
		lastMessage *discordgo.Message
	)

	var f = func(source mautrix.EventSource, evt *event.Event) {
		age := eventAge(evt.Timestamp)
		reportChannel, _ := channels[string(evt.RoomID)]

		if age < 60 && reportChannel != "" {
			sender := string(evt.Sender)

			lock.Lock()

			if lastSender != sender {
				sent = 0
				lastSender = sender
			}

			sent++

			logInfo.Printf("%s n%d %ds.", sender, sent, age)

			if sent == 1 {
				lastMessage, _ = session.ChannelMessageSend(
					reportChannel,
					fmt.Sprintf("**%s** sent a message in **Matrix**", sender),
				)

				lock.Unlock()
			} else if lastMessage != nil {
				lock.Unlock()

				session.ChannelMessageEdit(
					lastMessage.ChannelID,
					lastMessage.ID,
					fmt.Sprintf("**%s** sent %d messages in **Matrix**", sender, sent),
				)
			} else {
				lock.Unlock()
			}
		}

		client.MarkRead(evt.RoomID, evt.ID)
	}

	syncer := client.Syncer.(*mautrix.DefaultSyncer)

	syncer.OnEventType(event.EventMessage, f)
	syncer.OnEventType(event.EventEncrypted, f)

	logErr.Println(client.Sync())
}

func eventAge(timestamp int64) int64 {
	return time.Now().Unix() - (timestamp / 1000)
}
