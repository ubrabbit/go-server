package server

import (
	"net"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
)

type ClientUnit struct {
	sync.Mutex
	Address string
	Parent  *ServerUnit

	objectID int64
	session  cellnet.Session
}

func NewTcpClient(ev cellnet.Event) *ClientUnit {
	obj := new(ClientUnit)
	obj.Parent = nil
	obj.Address = ev.Session().Raw().(net.Conn).RemoteAddr().String()
	obj.session = ev.Session()
	obj.objectID = newObjectID()
	return obj
}

func (self *ClientUnit) ObjectID() int64 {
	return self.objectID
}

func (self *ClientUnit) Session() cellnet.Session {
	return self.session
}

func (self *ClientUnit) SessionID() int64 {
	return self.session.ID()
}

func (self *ClientUnit) PacketSend(msg interface{}) {
	self.Lock()
	defer self.Unlock()

	self.Session().Send(&msg)
}

func (self *ClientUnit) OnDisconnect() {

}
