package cellnet

import (
	"fmt"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"

	_ "github.com/davyxu/cellnet/peer/udp"
	_ "github.com/davyxu/cellnet/proc/udp"
)

import (
	. "github.com/ubrabbit/go-public/common"
)

type UdpClientHandle interface {
	OnProtoCommand(*UdpClient, interface{})
}

type UdpServer struct {
	sync.Mutex
	Name     string
	Address  string
	queueIns cellnet.EventQueue
	peerIns  cellnet.GenericPeer

	objectID     int64
	clientHandle interface{}
	waitStopped  chan bool
}

func NewUdpServer(name string, address string, handle interface{}) *UdpServer {
	obj := new(UdpServer)
	obj.Name = name
	obj.Address = address
	obj.clientHandle = handle
	obj.objectID = newObjectID()
	obj.waitStopped = make(chan bool, 1)

	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	queue := cellnet.NewEventQueue()
	peerIns := peer.NewGenericPeer("udp.Acceptor", name, address, queue)
	proc.BindProcessorHandler(peerIns, "udp.ltv", obj.packetRecv)
	obj.queueIns = queue
	obj.peerIns = peerIns

	go obj.serverRun()
	return obj
}

func (self *UdpServer) String() string {
	return fmt.Sprintf("[Server][%s][%s]-%d ", self.Address, self.Name, self.objectID)
}

func (self *UdpServer) WaitStop() {
	<-self.waitStopped
	LogInfo("Stopped: %v", self)
}

func (self *UdpServer) setStop() {
	if self.waitStopped != nil {
		self.waitStopped <- true
		self.waitStopped = nil
	}
}

//此函数运行失败就直接让它崩溃
func (self *UdpServer) serverRun() {
	// 开始侦听
	self.peerIns.Start()
	// 事件队列开始循环
	self.queueIns.StartLoop()
	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	self.queueIns.Wait()
}

func (self *UdpServer) packetRecv(ev cellnet.Event) {
	msg := ev.Message()
	LogDebug("packetRecv:  %d %v", ev.Session().ID(), msg)
	client := newUdpClient(ev)
	self.clientHandle.(UdpClientHandle).OnProtoCommand(client, msg)
}
