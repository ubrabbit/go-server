package server

import (
	"fmt"
	"net"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
)

import (
	. "github.com/ubrabbit/go-server/common"
)

type ClientUnit struct {
	sync.Mutex
	Address string
	Parent  *ServerUnit

	objectID  int64
	sessionID int64
	session   cellnet.Session
}

func NewTcpClient(ev cellnet.Event) *ClientUnit {
	obj := new(ClientUnit)
	obj.Parent = nil
	obj.Address = ev.Session().Raw().(net.Conn).RemoteAddr().String()
	obj.session = ev.Session()
	obj.sessionID = ev.Session().ID()
	obj.objectID = newObjectID()
	return obj
}

func (self *ClientUnit) String() string {
	return fmt.Sprintf("[Client][%s]-%d-%d ", self.Address, self.objectID, self.sessionID)
}

func (self *ClientUnit) ObjectID() int64 {
	return self.objectID
}

func (self *ClientUnit) Session() cellnet.Session {
	return self.session
}

func (self *ClientUnit) SessionID() int64 {
	return self.sessionID
}

func (self *ClientUnit) PacketSend(msg interface{}) {
	self.Lock()
	defer self.Unlock()

	self.Session().Send(&msg)
}

func (self *ClientUnit) OnConnectSucc() {
	LogInfo(self, "Connected")
}

func (self *ClientUnit) OnDisconnect() {
	LogInfo(self, "Disconnected")
}
