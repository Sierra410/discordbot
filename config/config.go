package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
)

type Config struct {
	lock sync.Mutex
	data map[string]interface{}
}

func (self *Config) Save(configFilePath string) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	m, err := json.MarshalIndent(self, "", "\t")
	if err != nil {
		return nil
	}

	err = ioutil.WriteFile(configFilePath, m, 0640)

	return err
}

func (self *Config) Load(configFilePath string) error {
	self.lock.Lock()
	defer self.lock.Unlock()

	bytes, err := ioutil.ReadFile(configFilePath)

	if os.IsNotExist(err) {
		self.Save(configFilePath)
		return nil
	} else if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &self)

	return err
}

func (self *Config) Get(key string, def interface{}) interface{} {
	self.lock.Lock()
	defer self.lock.Unlock()

	v, ok := self.data[key]
	if !ok {
		return def
	}

	return v
}
