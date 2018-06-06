package server

import (
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

	objectID  int64
	onCommand func(*ConnectUnit, interface{})
}

func NewTcpConnect(name string, address string) *ConnectUnit {
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
	obj.onCommand = nil

	proc.BindProcessorHandler(p, "tcp.ltv", obj.PacketRecv)
	obj.Run()
	return obj
}

func (self *ConnectUnit) ObjectID() int64 {
	return self.objectID
}

func (self *ConnectUnit) Session() cellnet.Session {
	return self.Peer.(interface {
		Session() cellnet.Session
	}).Session()
}

func (self *ConnectUnit) SessionID() int64 {
	return self.Session().ID()
}

func (self *ConnectUnit) SetCommand(f func(*ConnectUnit, interface{})) {
	self.Lock()
	defer self.Unlock()

	self.onCommand = f
}

func (self *ConnectUnit) Run() {
	// 开始发起到服务器的连接
	self.Peer.Start()
	// 事件队列开始循环
	self.Queue.StartLoop()
}

func (self *ConnectUnit) Disconnect() {
	self.Peer.Stop()
}

func (self *ConnectUnit) OnConnectSucc(ev cellnet.Event) {
	self.Lock()
	defer self.Unlock()

	pool := GetConnectPool()
	pool.Add(self)
}

func (self *ConnectUnit) OnDisconnect(ev cellnet.Event) {
	self.Lock()
	defer self.Unlock()

	pool := GetConnectPool()
	pool.Remove(self.SessionID())
}

func (self *ConnectUnit) PacketSend(msg interface{}) {
	self.Lock()
	defer self.Unlock()

	self.Session().Send(msg)
}

func (self *ConnectUnit) PacketRecv(ev cellnet.Event) {
	msg := ev.Message()
	switch msg.(type) {
	case *cellnet.SessionConnected:
		LogInfo("Connect connected")
		self.OnConnectSucc(ev)
	case *cellnet.SessionClosed:
		LogInfo("Connect error")
		self.OnDisconnect(ev)
	default:
		if self.onCommand != nil {
			self.onCommand(self, msg)
		} else {
			onConnectCommand(self, msg)
		}
	}
}