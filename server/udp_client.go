package server

import (
	"fmt"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
)

type UdpClient struct {
	sync.Mutex

	objectID int64
	session  cellnet.Session
}

func newUdpClient(ev cellnet.Event) *UdpClient {
	obj := new(UdpClient)
	obj.session = ev.Session()
	obj.objectID = newObjectID()
	return obj
}

func (self *UdpClient) String() string {
	return fmt.Sprintf("[UdpClient]-%d", self.objectID)
}

func (self *UdpClient) ObjectID() int64 {
	return self.objectID
}

func (self *UdpClient) Session() cellnet.Session {
	return self.session
}

func (self *UdpClient) PacketSend(msg interface{}) {
	self.Lock()
	defer self.Unlock()

	self.Session().Send(msg)
}
