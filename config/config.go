package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type Config struct {
	lock sync.RWMutex
	data map[string]interface{}

	path          string
	timer         *time.Timer
	timerDuration time.Duration
}

func NewConfig(path string) Config {
	return Config{
		lock: sync.RWMutex{},
		data: map[string]interface{}{},

		path:  path,
		timer: nil,
	}
}

func (self *Config) SetAutosaveTime(t time.Duration) { // t is in ms
	if t == 0 {
		self.timerDuration = 0
		self.timer.Stop()
		self.timer = nil
	} else {
		self.timerDuration = t

		if self.timer == nil {
			self.timer = time.AfterFunc(
				self.timerDuration,
				func() {
					err := self.Save()
					if err != nil {
						panic(err)
					}
				},
			)
		}
	}
}

func (self *Config) Save() error {
	self.lock.RLock()
	defer self.lock.RUnlock()

	m, err := json.MarshalIndent(self.data, "", "\t")
	if err != nil {
		return nil
	}

	err = ioutil.WriteFile(self.path, m, 0640)

	return err
}

func (self *Config) Load() error {
	self.lock.Lock()
	defer self.lock.Unlock()

	bytes, err := ioutil.ReadFile(self.path)

	if os.IsNotExist(err) {
		return self.Save()
	} else if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &self.data)
	if err != nil {
		return err
	}

	return nil
}

func (self *Config) Get(key string) interface{} {
	self.lock.RLock()
	defer self.lock.RUnlock()

	v, _ := self.data[key]

	return v
}

func (self *Config) UnsafeGet(key string) interface{} {
	v, _ := self.data[key]

	return v
}

func (self *Config) GetOr(key string, def interface{}) interface{} {
	self.lock.RLock()
	defer self.lock.RUnlock()

	v, ok := self.data[key]
	if !ok {
		return def
	}

	return v
}

func (self *Config) Set(key string, val interface{}) interface{} {
	self.lock.Lock()
	defer self.lock.Unlock()

	current, _ := self.data[key]

	self.data[key] = val

	if self.timer != nil {
		self.timer.Reset(self.timerDuration)
	}

	return current
}

func (self *Config) UnsafeSet(key string, val interface{}) interface{} {
	current, _ := self.data[key]

	self.data[key] = val

	if self.timer != nil {
		self.timer.Reset(self.timerDuration)
	}

	return current
}

func (self *Config) Lock() {
	self.lock.Lock()
}

func (self *Config) Unlock() {
	self.lock.Unlock()
}

func (self *Config) RLock() {
	self.lock.RLock()
}

func (self *Config) RUnlock() {
	self.lock.RUnlock()
}
