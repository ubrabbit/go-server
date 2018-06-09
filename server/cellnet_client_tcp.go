package server

import (
	"fmt"
	"net"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
)

type TcpClient struct {
	sync.Mutex
	Address string

	objectID  int64
	sessionID int64
	session   cellnet.Session
}

func NewTcpClient(ev cellnet.Event) *TcpClient {
	obj := new(TcpClient)
	obj.Address = ev.Session().Raw().(net.Conn).RemoteAddr().String()
	obj.session = ev.Session()
	obj.sessionID = ev.Session().ID()
	obj.objectID = newObjectID()
	return obj
}

func (self *TcpClient) String() string {
	return fmt.Sprintf("[TcpClient][%s]-%d-%d ", self.Address, self.objectID, self.sessionID)
}

func (self *TcpClient) ObjectID() int64 {
	return self.objectID
}

func (self *TcpClient) Session() cellnet.Session {
	return self.session
}

func (self *TcpClient) SessionID() int64 {
	return self.sessionID
}

func (self *TcpClient) PacketSend(msg interface{}) {
	self.Lock()
	defer self.Unlock()

	self.Session().Send(&msg)
}
