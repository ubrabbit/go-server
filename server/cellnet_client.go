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

type ClientUnit struct {
	sync.Mutex
	Name    string
	Address string
	Queue   cellnet.EventQueue
	Peer    cellnet.GenericPeer
}

func NewTcpClient(name string, address string) *ClientUnit {
	// 创建一个事件处理队列，整个客户端只有这一个队列处理事件，客户端属于单线程模型
	queue := cellnet.NewEventQueue()
	// 创建一个tcp的连接器，名称为client，连接地址为127.0.0.1:8801，将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	//p := peer.NewGenericPeer("tcp.Connector", "client", "127.0.0.1:18801", queue)
	p := peer.NewGenericPeer("tcp.Connector", name, address, queue)
	proc.BindProcessorHandler(p, "tcp.ltv", onClientCommand)
	//p.SetReconnectDuration(1)

	obj := new(ClientUnit)
	obj.Name = name
	obj.Address = address
	obj.Queue = queue
	obj.Peer = p
	pool := GetClientPool()
	pool.AddClient(obj)

	obj.Run()
	return obj
}

func (self *ClientUnit) Run() {
	// 开始发起到服务器的连接
	self.Peer.Start()
	// 事件队列开始循环
	self.Queue.StartLoop()
}

func (self *ClientUnit) Disconnect() {
	self.Peer.Stop()
}

func (self *ClientUnit) OnConnectSucc(ev cellnet.Event) {

}

func (self *ClientUnit) OnDisconnect(ev cellnet.Event) {

}

func (self *ClientUnit) PacketSend(msg interface{}) {
	self.Peer.(interface {
		Session() cellnet.Session
	}).Session().Send(msg)
}

func (self *ClientUnit) PacketRecv(ev cellnet.Event) {
	switch ev.Message().(type) {
	case *cellnet.SessionConnected:
		LogInfo("client connected")
		self.OnConnectSucc(ev)
	case *cellnet.SessionClosed:
		LogInfo("client error")
		self.OnDisconnect(ev)
	default:
		onClientCommand(ev)
	}
}
