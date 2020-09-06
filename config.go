package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

var configFilePath = "config.json"
var cfg = &Config{}

type Config struct {
	lock sync.Mutex

	Token                string
	Status               string
	DefaultCommandPrefix string
	Admins               []string

	PerServerConfig map[string]*PerServerConfig
}

func (self *Config) fillBlanks() {
	self.lock.Lock()

	if self.DefaultCommandPrefix == "" {
		self.DefaultCommandPrefix = "!"
	}

	if self.Admins == nil {
		self.Admins = []string{}
	}

	if self.PerServerConfig == nil {
		self.PerServerConfig = map[string]*PerServerConfig{}
	}

	self.lock.Unlock()
}

func (self *Config) GetPerServerConfig(id string) *PerServerConfig {
	self.lock.Lock()
	defer self.lock.Unlock()

	c, ok := self.PerServerConfig[id]
	if !ok {
		c = &PerServerConfig{}
		self.PerServerConfig[id] = c
	}

	return c
}

func (self *Config) Save() error {
	self.lock.Lock()
	defer self.lock.Unlock()

	self.fillBlanks()
	m, err := json.MarshalIndent(cfg, "", "\t")
	if err != nil {
		return nil
	}

	err = ioutil.WriteFile(configFilePath, m, 0640)

	return err
}

func (self *Config) Load() error {
	self.lock.Lock()
	defer self.lock.Unlock()

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
	self.lock.Lock()
	defer self.lock.Unlock()

	for _, adminId := range cfg.Admins {
		if s == adminId {
			return true
		}
	}

	return false
}
