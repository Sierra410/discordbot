package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

var configFilePath = "config.json"
var cfg = &Config{}

type Config struct {
	Token  string
	Status string

	CommandPrefix string

	Admins []string

	ReactionRoles map[string]ReactionRole
}

func (self *Config) fillBlanks() {
	// Defaults, basically

	if self.CommandPrefix == "" {
		self.CommandPrefix = "!"
	}

	if self.Admins == nil {
		self.Admins = []string{}
	}

	if self.ReactionRoles == nil {
		self.ReactionRoles = map[string]ReactionRole{}
	}
}

func (self *Config) Save() error {
	self.fillBlanks()
	m, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return nil
	}

	err = ioutil.WriteFile(configFilePath, m, 0640)

	return err
}

func (self *Config) Load() error {
	bytes, err := ioutil.ReadFile(configFilePath)

	if os.IsNotExist(err) {
		self.Save()
		return nil
	} else if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &cfg)
	cfg.fillBlanks()

	return err
}

func (self *Config) IsAdmin(s string) bool {
	for _, adminId := range cfg.Admins {
		if s == adminId {
			return true
		}
	}

	return false
}
