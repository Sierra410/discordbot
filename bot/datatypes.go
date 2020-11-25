package bot

import (
	"sync"
	"sync/atomic"
)

//
// Atomic Flag (bool)
//

type Flag struct {
	i int32
}

func NewFlag(b bool) *Flag {
	flag := &Flag{0}
	if b {
		flag.i = 1
	}

	return flag
}

func (self *Flag) Get() bool {
	return atomic.LoadInt32(&self.i) != 0
}

func (self *Flag) Set(b bool) {
	var i int32
	if b {
		i = 1
	}
	atomic.StoreInt32(&self.i, i)
}

//
// Thread Safe Map
//

type Smap struct {
	l *sync.Mutex
	m map[string]interface{}
}

func NewSmap() *Smap {
	return &Smap{
		&sync.Mutex{},
		map[string]interface{}{},
	}
}

func (self *Smap) Get(key string) (v interface{}, ok bool) {
	self.l.Lock()
	v, ok = self.m[key]
	self.l.Unlock()
	return v, ok
}

func (self *Smap) Set(key string, val interface{}) {
	self.l.Lock()
	self.m[key] = val
	self.l.Unlock()
}

func (self *Smap) Delete(key string) {
	self.l.Lock()
	delete(self.m, key)
	self.l.Unlock()
}
