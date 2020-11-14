package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

type ConfigInterface interface {
	Get(string) interface{}
	GetOr(string, interface{}) interface{}
	Set(string, interface{}) interface{}
	Save() error
	Load() error
}

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

	return self.unsafeSave()
}

func (self *Config) unsafeSave() error {
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

	return self.unsafeLoad()
}

func (self *Config) unsafeLoad() error {
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

	return self.unsafeGet
}

func (self *Config) unsafeGet(key string) interface{} {
	v, _ := self.data[key]

	return v
}

func (self *Config) GetOr(key string, def interface{}) interface{} {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.unsafeGetOr(key, def)
}

func (self *Config) unsafeGetOr(key string, def interface{}) interface{} {
	v, ok := self.data[key]
	if !ok {
		return def
	}

	return v
}

func (self *Config) Set(key string, val interface{}) interface{} {
	self.lock.Lock()
	defer self.lock.Unlock()

	return self.unsafeSet(key, val)
}

func (self *Config) unsafeSet(key string, val interface{}) interface{} {
	current, _ := self.data[key]

	self.data[key] = val

	if self.timer != nil {
		self.timer.Reset(self.timerDuration)
	}

	return current
}

func (self *Config) Lock(writing bool) *Handle {
	handle := &Handle{
		lock:    sync.RWMutex{},
		writing: writing,
		cfg:     self,
	}

	if writing {
		self.lock.Lock()
	} else {
		self.lock.RLock()
	}

	return handle
}

type Handle struct {
	lock    sync.RWMutex
	writing bool
	cfg     *Config
}

func (self *Handle) Drop() {
	self.lock.Lock()
	defer self.lock.Unlock()

	if self.writing {
		self.cfg.lock.Unlock()
	} else {
		self.cfg.lock.RUnlock()
	}

	self.cfg = nil
}

func (self *Handle) Get(key string) interface{} {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.cfg.unsafeGet(key)
}

func (self *Handle) GetOr(key string, def interface{}) interface{} {
	self.lock.RLock()
	defer self.lock.RUnlock()

	return self.cfg.unsafeGetOr(key, def)
}

func (self *Handle) Set(key string, val interface{}) interface{} {
	self.lock.Lock()
	defer self.lock.Unlock()

	return self.cfg.unsafeSet(key, val)
}

func (self *Handle) Save() error {
	self.lock.Lock()
	defer self.lock.Unlock()

	return self.cfg.unsafeSave()
}

func (self *Handle) Load() error {
	self.lock.Lock()
	defer self.lock.Unlock()

	return self.cfg.unsafeLoad()
}
