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

type ConfigInterface interface {
	Get(guildId string, key string) interface{}
	GetOr(guildId string, key string, value interface{}) interface{}
	Set(guildId string, key string, value interface{}) interface{}
	Save() error
	Load() error
}

var configFilePath = "config.json"
var cfg Config

// Extended behavior;

func toConfigKey(guildId, key string) string {
	if guildId == "" {
		return key
	}

	return guildId + "_" + key
}

func getPermissionLevel(config ConfigInterface, session *discordgo.Session, guildId string, id string) botPermissionLevel {
	botOwner := cfg.GetOr("", configBotOwner, "").(string)
	if botOwner != "" && id == botOwner {
		return botPermBotOwner
	}

	admins := interfaceToInterfaceSlice(config.Get("", configAdmins))
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

		admins = interfaceToInterfaceSlice(config.Get(guildId, configAdmins))
		if isInSlice(id, admins) {
			return botPermServerAdmin
		}
	}

	return botPermNone
}

func isAdmin(config ConfigInterface, session *discordgo.Session, guildId string, id string) bool {
	return getPermissionLevel(config, session, guildId, id) != botPermNone
}

func setPrefix(config ConfigInterface, guildId string, prefix string) string {
	return config.Set(guildId, configCommandPrefix, prefix).(string)
}

func getPrefix(config ConfigInterface, guildId string) string {
	// No prefix in DM
	if guildId == "" {
		return ""
	}

	p := interfaceToString(
		config.Get(guildId, configCommandPrefix),
	)

	if p != "" {
		return p
	}

	p = interfaceToString(
		config.Get("", configCommandPrefix),
	)

	if p != "" {
		return p
	}

	config.Set("", configCommandPrefix, configDefaultCommandPrefix)

	return configDefaultCommandPrefix
}

// Config is a simple key:value storage with several Discordgo specific helper funciton
type Config struct {
	cfg config.Config
}

func NewConfig(path string) Config {
	return Config{
		cfg: config.NewConfig(path),
	}
}

// Extension

func (self *Config) IsAdmin(session *discordgo.Session, guildId string, id string) bool {
	return isAdmin(self, session, guildId, id)
}

func (self *Config) GetPermissionLevel(session *discordgo.Session, guildId string, id string) botPermissionLevel {
	return getPermissionLevel(self, session, guildId, id)
}

func (self *Config) SetPrefix(guildId string, prefix string) string {
	return setPrefix(self, guildId, prefix)
}

func (self *Config) GetPrefix(guildId string) string {
	return getPrefix(self, guildId)
}

// DefaultsConfig

func (self *Config) Set(guildId string, k string, v interface{}) interface{} {
	return self.cfg.Set(toConfigKey(guildId, k), v)
}

func (self *Config) Get(guildId string, k string) interface{} {
	return self.cfg.Get(toConfigKey(guildId, k))
}

func (self *Config) GetOr(guildId string, k string, def interface{}) interface{} {
	return self.cfg.GetOr(toConfigKey(guildId, k), def)
}

func (self *Config) Save() error {
	return self.cfg.Save()
}

func (self *Config) Load() error {
	return self.cfg.Load()
}

func (self *Config) Lock(writing bool) *ConfigHandle {
	return &ConfigHandle{
		handle: self.cfg.Lock(writing),
	}
}

type ConfigHandle struct {
	handle *config.Handle
}

// Extension

func (self *ConfigHandle) IsAdmin(session *discordgo.Session, guildId string, id string) bool {
	return isAdmin(self, session, guildId, id)
}

func (self *ConfigHandle) GetPermissionLevel(session *discordgo.Session, guildId string, id string) botPermissionLevel {
	return getPermissionLevel(self, session, guildId, id)
}

func (self *ConfigHandle) SetPrefix(guildId string, prefix string) string {
	return setPrefix(self, guildId, prefix)
}

func (self *ConfigHandle) GetPrefix(guildId string) string {
	return getPrefix(self, guildId)
}

// Defaults

func (self *ConfigHandle) Drop() {
	self.handle.Drop()
}

func (self *ConfigHandle) Set(g, k string, v interface{}) interface{} {
	return self.handle.Set(toConfigKey(g, k), v)
}

func (self *ConfigHandle) Get(g, k string) interface{} {
	return self.handle.Get(toConfigKey(g, k))
}

func (self *ConfigHandle) GetOr(g, k string, def interface{}) interface{} {
	return self.handle.GetOr(toConfigKey(g, k), def)
}

func (self *ConfigHandle) Save() error {
	return self.handle.Save()
}

func (self *ConfigHandle) Load() error {
	return self.handle.Load()
}
