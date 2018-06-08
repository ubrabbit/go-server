package server

import (
	"fmt"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

import (
	. "github.com/ubrabbit/go-server/common"
)

type ConnectUnit struct {
	sync.Mutex
	Name    string
	Address string
	Queue   cellnet.EventQueue
	Peer    cellnet.GenericPeer

	sessionID     int64
	objectID      int64
	waitConnected chan bool
	onCommand     func(*ConnectUnit, interface{})
	eventTrigger  func(string, ...interface{})
}

func NewTcpConnect(name string, address string, f1 func(*ConnectUnit, interface{}), f2 func(string, ...interface{})) *ConnectUnit {
	// 创建一个事件处理队列，整个客户端只有这一个队列处理事件，客户端属于单线程模型
	queue := cellnet.NewEventQueue()
	// 创建一个tcp的连接器，名称为Connect，连接地址为127.0.0.1:8801，将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	//p := peer.NewGenericPeer("tcp.Connector", "Connect", "127.0.0.1:18801", queue)
	p := peer.NewGenericPeer("tcp.Connector", name, address, queue)
	//p.SetReconnectDuration(1)

	obj := new(ConnectUnit)
	obj.Name = name
	obj.Address = address
	obj.Queue = queue
	obj.Peer = p
	obj.objectID = newObjectID()
	obj.sessionID = 0
	obj.onCommand = f1
	obj.eventTrigger = f2
	obj.waitConnected = make(chan bool, 1)

	proc.BindProcessorHandler(p, "tcp.ltv", obj.PacketRecv)
	// 开始发起到服务器的连接
	obj.Peer.Start()
	// 事件队列开始循环
	obj.Queue.StartLoop()

	//等待连接成功再返回
	<-obj.waitConnected
	obj.waitConnected = nil
	return obj
}

//__repr__
func (self *ConnectUnit) String() string {
	return fmt.Sprintf("[Connect][%s]-%d-%d ", self.Address, self.objectID, self.sessionID)
}

func (self *ConnectUnit) ObjectID() int64 {
	return self.objectID
}

func (self *ConnectUnit) Session() cellnet.Session {
	return self.Peer.(interface {
		Session() cellnet.Session
	}).Session()
}

//sessionID在断线后通过Session获取不到
func (self *ConnectUnit) SessionID() int64 {
	return self.sessionID
}

func (self *ConnectUnit) Disconnect() {
	self.Lock()
	defer func() {
		err := recover()
		if err != nil {
			LogError(self, " Disconnect Error: ", err)
		}
		self.Unlock()
	}()
	self.Peer.Stop()
}

func (self *ConnectUnit) OnConnectSucc(ev cellnet.Event) {
	self.Lock()
	defer self.Unlock()

	LogInfo(self, "ConnectSucc")
	self.sessionID = self.Session().ID()
	pool := GetConnectPool()
	pool.Add(self)

	//连接成功，取消阻塞
	self.waitConnected <- true
	if self.eventTrigger != nil {
		self.eventTrigger("Connect", self.Address)
	}
}

func (self *ConnectUnit) OnDisconnect(ev cellnet.Event) {
	self.Lock()
	defer self.Unlock()

	LogInfo(self, "Disconnected")
	pool := GetConnectPool()
	pool.Remove(self.SessionID())

	if self.eventTrigger != nil {
		self.eventTrigger("DisConnect", self.Address)
	}
}

func (self *ConnectUnit) PacketSend(msg interface{}) {
	self.Lock()
	defer func() {
		err := recover()
		if err != nil {
			LogError(self, "PacketSend: ", err)
		}
		self.Unlock()
	}()
	self.Session().Send(msg)
}

func (self *ConnectUnit) PacketRecv(ev cellnet.Event) {
	defer func() {
		err := recover()
		if err != nil {
			LogError(self, "PacketRecv: ", err)
		}
	}()

	msg := ev.Message()
	switch msg.(type) {
	case *cellnet.SessionConnected:
		self.OnConnectSucc(ev)
	case *cellnet.SessionClosed:
		self.OnDisconnect(ev)
	default:
		self.onCommand(self, msg)
	}
}
