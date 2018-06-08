package server

import (
	"fmt"
	"net"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
)

type Client struct {
	sync.Mutex
	Address string
	Parent  *ServerUnit

	objectID  int64
	sessionID int64
	session   cellnet.Session
}

func NewTcpClient(ev cellnet.Event) *Client {
	obj := new(Client)
	obj.Parent = nil
	obj.Address = ev.Session().Raw().(net.Conn).RemoteAddr().String()
	obj.session = ev.Session()
	obj.sessionID = ev.Session().ID()
	obj.objectID = newObjectID()
	return obj
}

func (self *Client) String() string {
	return fmt.Sprintf("[Client][%s]-%d-%d ", self.Address, self.objectID, self.sessionID)
}

func (self *Client) ObjectID() int64 {
	return self.objectID
}

func (self *Client) Session() cellnet.Session {
	return self.session
}

func (self *Client) SessionID() int64 {
	return self.sessionID
}

func (self *Client) PacketSend(msg interface{}) {
	self.Lock()
	defer self.Unlock()

	self.Session().Send(&msg)
}
