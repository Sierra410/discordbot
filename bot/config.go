package bot

import (
	"local/discordbot/config"

	"github.com/bwmarrin/discordgo"
)

const (
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
	self.RLock()
	defer self.RUnlock()

	admins := interfaceToInterfaceSlice(self.Get("", configAdmins))
	if isInSlice(id, admins) {
		return true
	}

	admins = interfaceToInterfaceSlice(self.Get(guildId, configAdmins))
	if isInSlice(id, admins) {
		return true
	}

	guild, err := session.Guild(guildId)
	return err == nil && guild.OwnerID == id
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
