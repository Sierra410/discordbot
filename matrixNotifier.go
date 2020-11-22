package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
)

func init() {
	addReadyCallback(matrixChatWatchdog)
}

type recentMessage struct {
	message *discordgo.Message
	sender  string
	sent    uint32
}

type recentMessages struct {
	lock              sync.Mutex
	timeout           time.Duration
	messagesByChannel map[string]*recentMessage
	cleanupDelay      map[string]*time.Timer
}

func newRecentMessages(timeoutSeconds time.Duration) *recentMessages {
	return &recentMessages{
		lock:              sync.Mutex{},
		timeout:           time.Second * timeoutSeconds,
		cleanupDelay:      map[string]*time.Timer{},
		messagesByChannel: map[string]*recentMessage{},
	}
}

func (self *recentMessages) getMostRecent(channelId string) *recentMessage {
	self.lock.Lock()
	defer self.lock.Unlock()

	rm, ok := self.messagesByChannel[channelId]
	if !ok {
		return nil
	}

	return rm
}

func (self *recentMessages) addNewMessage(message *discordgo.Message, sender string, sent uint32) {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.messagesByChannel[message.ChannelID] = &recentMessage{
		sender:  sender,
		message: message,
		sent:    sent,
	}

	timer, ok := self.cleanupDelay[message.ChannelID]
	if !ok {
		self.cleanupDelay[message.ChannelID] = time.AfterFunc(
			self.timeout,
			func() {
				self.lock.Lock()
				defer self.lock.Unlock()

				self.messagesByChannel[message.ChannelID] = nil
			},
		)
	} else {
		timer.Reset(self.timeout)
	}
}

func (self *recentMessages) cleanUp() {
	self.lock.Lock()
	defer self.lock.Unlock()

	for k, v := range self.messagesByChannel {
		if v == nil {
			continue
		}

		timeStamp, err := discordgo.SnowflakeTimestamp(v.message.ID)
		if err != nil || time.Since(timeStamp) > self.timeout {
			self.messagesByChannel[k] = nil
		}
	}
}

func matrixChatWatchdog(session *discordgo.Session, ready *discordgo.Ready) {
	prefix := "matrix"
	homeserver := cfg.GetOr(prefix, "homeserver", "matrix.org").(string)
	username := cfg.Get(prefix, "username").(string)
	password := cfg.Get(prefix, "password").(string)
	channels := map[string][]string{}

	for k, v := range cfg.Get(prefix, "channels").(map[string]interface{}) {
		channels[k] = interfaceToStringSlice(v)
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

	recentMessages := newRecentMessages(60)

	var f = func(source mautrix.EventSource, evt *event.Event) {
		defer client.MarkRead(evt.RoomID, evt.ID)

		age := eventAge(evt.Timestamp)
		if age > 60 {
			return
		}

		sender := string(evt.Sender)
		reportChannels, _ := channels[string(evt.RoomID)]

		logInfo.Printf("From %s to %v\n", sender, reportChannels)

		for _, reportChannel := range reportChannels {
			rm := recentMessages.getMostRecent(reportChannel)
			if rm == nil || rm.sender != sender {
				message, _ := session.ChannelMessageSend(
					reportChannel,
					fmt.Sprintf("**%s** sent a message in **Matrix**", sender),
				)

				recentMessages.addNewMessage(message, sender, 1)
			} else {
				atomic.AddUint32(&rm.sent, 1)

				session.ChannelMessageEdit(
					rm.message.ChannelID,
					rm.message.ID,
					fmt.Sprintf("**%s** sent %d messages in **Matrix**", sender, rm.sent),
				)
			}
		}
	}

	syncer := client.Syncer.(*mautrix.DefaultSyncer)

	syncer.OnEventType(event.EventMessage, f)
	syncer.OnEventType(event.EventEncrypted, f)

	logErr.Println(client.Sync())
}

func eventAge(timestamp int64) int64 {
	return time.Now().Unix() - (timestamp / 1000)
}
