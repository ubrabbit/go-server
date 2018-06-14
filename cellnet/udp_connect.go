package cellnet

import (
	"fmt"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/udp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/udp"
)

import (
	. "github.com/ubrabbit/go-common/common"
)

type UdpConnectHandle interface {
	OnProtoCommand(*UdpConnect, interface{})
}

type UdpConnect struct {
	sync.Mutex
	Name     string
	Address  string
	queueIns cellnet.EventQueue
	peerIns  cellnet.GenericPeer

	objectID      int64
	waitConnected chan bool

	connectHandle interface{}
}

func NewUdpConnect(name string, address string, handle interface{}) *UdpConnect {
	obj := new(UdpConnect)
	obj.Name = name
	obj.Address = address
	obj.objectID = newObjectID()
	obj.connectHandle = handle
	obj.waitConnected = make(chan bool, 1)

	// 创建一个事件处理队列，整个客户端只有这一个队列处理事件，客户端属于单线程模型
	queue := cellnet.NewEventQueue()
	// 创建一个tcp的连接器，名称为Connect，连接地址为127.0.0.1:8801，将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	// peer.NewGenericPeer("tcp.Connector", "Connect", "127.0.0.1:18801", queue)
	peerIns := peer.NewGenericPeer("udp.Connector", name, address, queue)
	proc.BindProcessorHandler(peerIns, "udp.ltv", obj.packetRecv)
	obj.queueIns = queue
	obj.peerIns = peerIns
	// 开始发起到服务器的连接
	obj.peerIns.Start()
	// 事件队列开始循环
	obj.queueIns.StartLoop()

	// 等待连接成功再返回
	<-obj.waitConnected
	obj.waitConnected = nil
	return obj
}

//__repr__
func (self *UdpConnect) String() string {
	return fmt.Sprintf("[UdpConnect][%s]-%d", self.Address, self.objectID)
}

func (self *UdpConnect) ObjectID() int64 {
	return self.objectID
}

func (self *UdpConnect) Session() cellnet.Session {
	return self.peerIns.(interface {
		Session() cellnet.Session
	}).Session()
}

func (self *UdpConnect) OnConnectSucc(ev cellnet.Event) {
	LogInfo(self, "ConnectSucc")

	//连接成功，取消阻塞
	if self.waitConnected != nil {
		self.waitConnected <- true
		self.waitConnected = nil
	}
}

func (self *UdpConnect) PacketSend(msg interface{}) {
	self.Lock()
	defer func() {
		err := recover()
		if err != nil {
			LogError(self, "PacketSend Error: ", err)
		}
		self.Unlock()
	}()
	self.Session().Send(msg)
}

func (self *UdpConnect) packetRecv(ev cellnet.Event) {
	defer func() {
		err := recover()
		if err != nil {
			LogError(self, "packetRecv Error: ", err)
		}
	}()
	LogInfo("packetRecv")
	msg := ev.Message()
	switch msg.(type) {
	case *cellnet.SessionConnected:
		self.OnConnectSucc(ev)
	default:
		self.connectHandle.(UdpConnectHandle).OnProtoCommand(self, msg)
	}
}
