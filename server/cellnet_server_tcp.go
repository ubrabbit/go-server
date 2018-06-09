package server

import (
	"fmt"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"

	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/cellnet/rpc"
)

import (
	. "github.com/ubrabbit/go-server/common"
)

type TcpClientHandle interface {
	OnProtoCommand(*TcpClient, interface{})
	OnRpcCommand(*TcpClient, interface{}) (interface{}, error)
	OnEventTrigger(*TcpClient, string, ...interface{})
}

type TcpServer struct {
	sync.Mutex
	Name    string
	Address string
	Queue   cellnet.EventQueue
	Peer    cellnet.GenericPeer
	Pool    map[int64]*TcpClient

	objectID     int64
	clientHandle interface{}
	waitStopped  chan bool
}

func NewTcpServer(name string, address string, handle interface{}) *TcpServer {
	obj := new(TcpServer)
	obj.Name = name
	obj.Address = address
	obj.Pool = make(map[int64]*TcpClient, 0)
	obj.clientHandle = handle
	obj.objectID = newObjectID()
	obj.waitStopped = make(chan bool, 1)

	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	queue := cellnet.NewEventQueue()
	// 创建一个tcp的侦听器，名称为server，连接地址为127.0.0.1:8801，所有连接将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	peerIns := peer.NewGenericPeer("tcp.Acceptor", name, address, queue)
	proc.BindProcessorHandler(peerIns, "tcp.ltv", obj.PacketRecv)
	obj.Queue = queue
	obj.Peer = peerIns

	go obj.serverRun()
	return obj
}

func (self *TcpServer) String() string {
	return fmt.Sprintf("[Server][%s][%s]-%d ", self.Address, self.Name, self.objectID)
}

func (self *TcpServer) WaitStop() {
	<-self.waitStopped
	LogInfo(self, "Stopped")
}

func (self *TcpServer) setStop() {
	if self.waitStopped != nil {
		self.waitStopped <- true
		self.waitStopped = nil
	}
}

//此函数运行失败就直接让它崩溃
func (self *TcpServer) serverRun() {
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

func (self *TcpServer) Disconnect() {
	self.Lock()
	defer func() {
		err := recover()
		if err != nil {
			LogError(self, " Disconnect Error: ", err)
		}
		self.setStop()
		self.Unlock()
	}()
	self.Peer.Stop()
}

func (self *TcpServer) OnConnectSucc(ev cellnet.Event) {
	self.Lock()
	defer self.Unlock()

	client := NewTcpClient(ev)
	self.Pool[client.SessionID()] = client

	LogInfo(client, "Connected")
	self.clientHandle.(ClientHandle).OnEventTrigger(client, "Connect")
}

func (self *TcpServer) OnDisconnect(ev cellnet.Event) {
	self.Lock()
	defer self.Unlock()

	sessionID := ev.Session().ID()
	client, ok := self.Pool[sessionID]
	if ok {
		delete(self.Pool, sessionID)
		LogInfo(client, "Disconnected")
		self.clientHandle.(ClientHandle).OnEventTrigger(client, "DisConnect")
	}
}

func (self *TcpServer) GetClient(sessionID int64) *TcpClient {
	self.Lock()
	defer self.Unlock()

	obj, ok := self.Pool[sessionID]
	if ok {
		return obj
	}
	return nil
}

func (self *TcpServer) PacketRecv(ev cellnet.Event) {
	//LogInfo("PacketRecv:  ", ev.Session().ID())
	msg := ev.Message()
	switch msg.(type) {
	// 有新的连接
	case *cellnet.SessionAccepted:
		LogInfo("Server Accepted", ev.Session().ID())
		self.OnConnectSucc(ev)
	// 有连接断开
	case *cellnet.SessionClosed:
		LogInfo("Session Closed: ", ev.Session().ID())
		self.OnDisconnect(ev)
	default:
		client := self.GetClient(ev.Session().ID())
		if client == nil {
			LogError(self, "Invalid Client: ", ev.Session().ID())
		} else {
			// 当服务器收到的是一个rpc消息
			if rpcevent, ok := ev.(*rpc.RecvMsgEvent); ok {
				response, err := self.clientHandle.(ClientHandle).OnRpcCommand(client, msg)
				if err != nil {
					LogError(self, "RpcCommand Error")
				} else {
					if response != nil {
						rpcevent.Reply(response)
					}
				}
			} else {
				self.clientHandle.(ClientHandle).OnProtoCommand(client, msg)
			}
		}
	}
}

func (self *TcpServer) Broadcast(msg interface{}) {
	self.Peer.(cellnet.SessionAccessor).VisitSession(
		func(ses cellnet.Session) bool {
			ses.Send(msg)
			return true
		})
}
