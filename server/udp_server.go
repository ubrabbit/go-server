package server

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
	. "github.com/ubrabbit/go-server/common"
)

type UdpClientHandle interface {
	OnProtoCommand(*UdpClient, interface{})
}

type UdpServer struct {
	sync.Mutex
	Name    string
	Address string
	Queue   cellnet.EventQueue
	Peer    cellnet.GenericPeer

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
	// 创建一个tcp的侦听器，名称为server，连接地址为127.0.0.1:8801，所有连接将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	peerIns := peer.NewGenericPeer("udp.Acceptor", name, address, queue)
	proc.BindProcessorHandler(peerIns, "udp.ltv", obj.PacketRecv)
	obj.Queue = queue
	obj.Peer = peerIns

	go obj.serverRun()
	return obj
}

func (self *UdpServer) String() string {
	return fmt.Sprintf("[Server][%s][%s]-%d ", self.Address, self.Name, self.objectID)
}

func (self *UdpServer) WaitStop() {
	<-self.waitStopped
	LogInfo(self, "Stopped")
}

func (self *UdpServer) setStop() {
	if self.waitStopped != nil {
		self.waitStopped <- true
		self.waitStopped = nil
	}
}

//此函数运行失败就直接让它崩溃
func (self *UdpServer) serverRun() {
	defer func() {
		err := recover()
		if err != nil {
			LogError(self, "RunError:  ", err)
		}
		self.Lock()
		self.setStop()
		self.Unlock()
	}()
	// 开始侦听
	self.Peer.Start()
	// 事件队列开始循环
	self.Queue.StartLoop()
	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	self.Queue.Wait()
}

func (self *UdpServer) PacketRecv(ev cellnet.Event) {
	LogInfo("PacketRecv:  ", ev.Session().ID())
	msg := ev.Message()
	LogInfo("111")
	client := newUdpClient(ev)
	LogInfo("222")
	self.clientHandle.(UdpClientHandle).OnProtoCommand(client, msg)
}
