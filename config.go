package main

import (
	"discordbot/config"

	"github.com/bwmarrin/discordgo"
)

const (
	configBotOwner      = "BotOwner"
	configStatus        = "Status"
	configAdmins        = "Admins"
	configCommandPrefix = "CommandPrefix"

	configDefaultCommandPrefix = "!"
)

var configFilePath = "config.json"
var cfg Config

type Config struct {
	cfg config.Config
}

func NewConfig(path string) Config {
	return Config{
		cfg: config.NewConfig(path),
	}
}

func (self *Config) IsAdmin(session *discordgo.Session, guildId string, id string) bool {
	return self.GetPermissionLevel(session, guildId, id) != botPermNone
}

func (self *Config) GetPermissionLevel(session *discordgo.Session, guildId string, id string) botPermissionLevel {
	self.cfg.Lock()
	defer self.cfg.Unlock()

	return self.UnsafeGetPermissionLevel(session, guildId, id)
}

func (self *Config) UnsafeGetPermissionLevel(session *discordgo.Session, guildId string, id string) botPermissionLevel {
	botOwner := cfg.UnsafeGetOr("", configBotOwner, "").(string)
	if botOwner != "" && id == botOwner {
		return botPermBotOwner
	}

	admins := interfaceToInterfaceSlice(self.UnsafeGet("", configAdmins))
	if isInSlice(id, admins) {
		return botPermBotAdmin
	}

	if guildId != "" {
		guild, err := session.Guild(guildId)
		if err != nil {
			panic(err)
		}

		if guild.OwnerID == id {
			return botPermServerOwner
		}

		admins = interfaceToInterfaceSlice(self.UnsafeGet(guildId, configAdmins))
		if isInSlice(id, admins) {
			return botPermServerAdmin
		}
	}

	return botPermNone
}

func (self *Config) SetPrefix(guildId string, prefix string) string {
	return self.Set(guildId, configCommandPrefix, prefix).(string)
}

func (self *Config) GetPrefix(guildId string) string {
	if guildId == "" {
		return ""
	}

	p := interfaceToString(
		self.Get(guildId, configCommandPrefix),
	)

	if p != "" {
		return p
	}

	p = interfaceToString(
		self.Get("", configCommandPrefix),
	)

	if p != "" {
		return p
	}

	cfg.Set("", configCommandPrefix, configDefaultCommandPrefix)

	return configDefaultCommandPrefix
}

func (self *Config) GetOr(guildId string, k string, def interface{}) interface{} {
	x := self.Get(guildId, k)
	if x == nil {
		return def
	}

	return x
}

func (self *Config) UnsafeGetOr(guildId string, k string, def interface{}) interface{} {
	x := self.UnsafeGet(guildId, k)
	if x == nil {
		return def
	}

	return x
}

func (self *Config) Get(guildId string, k string) interface{} {
	if guildId != "" {
		k = guildId + "_" + k
	}

	return self.cfg.Get(k)
}

func (self *Config) UnsafeGet(guildId string, k string) interface{} {
	if guildId != "" {
		k = guildId + "_" + k
	}

	return self.cfg.UnsafeGet(k)
}

func (self *Config) Set(guildId string, k string, v interface{}) interface{} {
	if guildId != "" {
		k = guildId + "_" + k
	}

	return self.cfg.Set(k, v)
}

func (self *Config) UnsafeSet(guildId string, k string, v interface{}) interface{} {
	if guildId != "" {
		k = guildId + "_" + k
	}

	return self.cfg.UnsafeSet(k, v)
}

func (self *Config) Save() error {
	return self.cfg.Save()
}

func (self *Config) Load() error {
	return self.cfg.Load()
}

func (self *Config) Lock() {
	self.cfg.Lock()
}

func (self *Config) Unlock() {
	self.cfg.Unlock()
}

func (self *Config) RLock() {
	self.cfg.RLock()
}

func (self *Config) RUnlock() {
	self.cfg.RUnlock()
}
