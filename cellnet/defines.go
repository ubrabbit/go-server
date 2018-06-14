package cellnet

import (
	"sync"
)

var (
	g_ObjectID *ObjectID = nil
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

func InitServer() {
}
