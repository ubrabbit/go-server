package server

import (
	"sync"
)

type ClientPool struct {
	sync.Mutex
	Pool map[string]*ClientUnit
}

var (
	g_ClientPool *ClientPool = nil
)

func InitClientPool() {
	g_ClientPool = new(ClientPool)
	g_ClientPool.Pool = make(map[string]*ClientUnit, 0)
}

func GetClientPool() *ClientPool {
	return g_ClientPool
}

func (self *ClientPool) AddClient(client *ClientUnit) {
	self.Lock()
	defer self.Unlock()
	self.Pool[client.Name] = client
}

func (self *ClientPool) RemoveClient(name string) bool {
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

func (self *ClientPool) GetClient(name string) *ClientUnit {
	self.Lock()
	defer self.Unlock()
	obj, ok := self.Pool[name]
	if ok {
		return obj
	}
	return nil
}
