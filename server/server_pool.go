package server

import (
	"sync"
)

type ServerPool struct {
	sync.Mutex
	Pool map[string]*ServerUnit
}

var (
	g_ServerPool *ServerPool = nil
)

func InitServerPool() {
	g_ServerPool = new(ServerPool)
	g_ServerPool.Pool = make(map[string]*ServerUnit, 0)
}

func GetServerPool() *ServerPool {
	return g_ServerPool
}

func (self *ServerPool) AddServer(obj *ServerUnit) {
	self.Lock()
	defer self.Unlock()
	self.Pool[obj.Name] = obj
}

func (self *ServerPool) RemoveServer(name string) bool {
	self.Lock()
	defer self.Unlock()

	ob, ok := self.Pool[name]
	if ok {
		delete(self.Pool, name)
		ob.Disconnect()
		return true
	}
	return false
}

func (self *ServerPool) GetServer(name string) *ServerUnit {
	self.Lock()
	defer self.Unlock()
	obj, ok := self.Pool[name]
	if ok {
		return obj
	}
	return nil
}
