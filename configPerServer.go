package main

func NewPerServerConfig(prefix string) *PerServerConfig {
	return &PerServerConfig{
		Admins:        []string{},
		CommandPrefix: prefix,
		data:          map[string]string{},
		hasChanged:    false,
	}
}

type PerServerConfig struct {
	Admins        []string
	CommandPrefix string
	data          map[string]string
	hasChanged    bool
}

func (self *PerServerConfig) Set(k, v string) {
	self.data[k] = v
	self.hasChanged = true
}

func (self *PerServerConfig) Get(k string) (string, bool) {
	v, ok := self.data[k]
	return v, ok
}
