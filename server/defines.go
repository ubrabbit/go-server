package server

import (
	"sync"
)

var (
	g_ObjectID    *ObjectID    = nil
	g_ServerPool  *ServerPool  = nil
	g_ConnectPool *ConnectPool = nil
)

type ObjectID struct {
	sync.Mutex
	UUID int64
}

func newObjectID() int64 {
	if g_ObjectID == nil {
		g_ObjectID = new(ObjectID)
		g_ObjectID.UUID = 10000
	}
	g_ObjectID.Lock()
	defer g_ObjectID.Unlock()

	g_ObjectID.UUID++
	return g_ObjectID.UUID
}

func InitConnectPool() {
	g_ConnectPool = new(ConnectPool)
	g_ConnectPool.Pool = make(map[int64]*ConnectUnit, 0)
}

func InitServerPool() {
	g_ServerPool = new(ServerPool)
	g_ServerPool.Pool = make(map[string]*ServerUnit, 0)
}

func GetServerPool() *ServerPool {
	return g_ServerPool
}

func GetConnectPool() *ConnectPool {
	return g_ConnectPool
}

func InitServer() {
	InitConnectPool()
	InitServerPool()
}
