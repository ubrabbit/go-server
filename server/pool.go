package server

import (
	"sync"
)

type ConnectPool struct {
	sync.Mutex
	Pool map[int64]*ConnectUnit
}

type ServerPool struct {
	sync.Mutex
	Pool map[string]*ServerUnit
}

func (self *ServerPool) Add(obj *ServerUnit) {
	self.Lock()
	defer self.Unlock()
	self.Pool[obj.Name] = obj
}

func (self *ServerPool) Remove(name string) bool {
	self.Lock()
	defer self.Unlock()

	_, ok := self.Pool[name]
	if ok {
		delete(self.Pool, name)
	}
	return ok
}

func (self *ServerPool) Get(name string) *ServerUnit {
	self.Lock()
	defer self.Unlock()
	obj, ok := self.Pool[name]
	if ok {
		return obj
	}
	return nil
}

func (self *ConnectPool) Add(c *ConnectUnit) {
	self.Lock()
	defer self.Unlock()
	self.Pool[c.SessionID()] = c
}

func (self *ConnectPool) Remove(sessionID int64) bool {
	self.Lock()
	defer self.Unlock()

	_, ok := self.Pool[sessionID]
	if ok {
		delete(self.Pool, sessionID)
	}
	return ok
}

func (self *ConnectPool) Get(sessionID int64) *ConnectUnit {
	self.Lock()
	defer self.Unlock()
	obj, ok := self.Pool[sessionID]
	if ok {
		return obj
	}
	return nil
}
