package cellnet

import (
	"fmt"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
)

import (
	. "github.com/ubrabbit/go-common/common"
)

type TcpClient struct {
	sync.Mutex
	Address string

	objectID  int64
	sessionID int64
	session   cellnet.Session
}

func newTcpClient(address string, ev cellnet.Event) *TcpClient {
	obj := new(TcpClient)
	//obj.Address = ev.Session().Raw().(net.Conn).RemoteAddr().String()
	obj.Address = address
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

	session := self.Session()
	if session == nil {
		LogError("Session Closed: ", self)
		return
	}
	session.Send(msg)
}
