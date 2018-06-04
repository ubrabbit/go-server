package server

import (
	"sync"
)

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"

	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
)

type ServerUnit struct {
	sync.Mutex
	Name    string
	Address string
	Queue   cellnet.EventQueue
	Peer    cellnet.GenericPeer
}

func NewTcpServer(name string, address string, is_block bool) {
	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	queue := cellnet.NewEventQueue()
	// 创建一个tcp的侦听器，名称为server，连接地址为127.0.0.1:8801，所有连接将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	p := peer.NewGenericPeer("tcp.Acceptor", name, address, queue)
	proc.BindProcessorHandler(p, "tcp.ltv", onServerCommand)

	obj := new(ServerUnit)
	obj.Name = name
	obj.Address = address
	obj.Queue = queue
	obj.Peer = p

	pool := GetServerPool()
	pool.AddServer(obj)

	if is_block {
		obj.Run()
	} else {
		go obj.Run()
	}
}

func (self *ServerUnit) Run() {
	// 开始侦听
	self.Peer.Start()
	// 事件队列开始循环
	self.Queue.StartLoop()
	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	self.Queue.Wait()
}

func (self *ServerUnit) Disconnect() {
	self.Peer.Stop()
}
