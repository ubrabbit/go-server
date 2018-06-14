package cellnet

import (
	"fmt"
	"net"
	"sync"
)

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	"github.com/davyxu/cellnet/proc"

	_ "github.com/davyxu/cellnet/peer/gorillaws"
	_ "github.com/davyxu/cellnet/peer/tcp"
	_ "github.com/davyxu/cellnet/proc/gorillaws"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/cellnet/rpc"
	"github.com/gorilla/websocket"
)

import (
	. "github.com/ubrabbit/go-common/common"
)

type TcpClientHandle interface {
	OnProtoCommand(*TcpClient, interface{})
	OnRpcCommand(*TcpClient, interface{}) (interface{}, error)
	OnEventTrigger(*TcpClient, string, ...interface{})
}

type TcpServer struct {
	sync.Mutex
	Name          string
	Type          string
	listenAddress string
	queueIns      cellnet.EventQueue
	peerIns       cellnet.GenericPeer
	clientPool    map[int64]*TcpClient

	objectID     int64
	clientHandle interface{}
	waitStopped  chan bool
}

func NewTcpServer(name string, address string, handle interface{}) *TcpServer {
	obj := new(TcpServer)
	obj.Name = name
	obj.Type = "tcp"
	obj.listenAddress = address
	obj.clientPool = make(map[int64]*TcpClient, 0)
	obj.clientHandle = handle
	obj.objectID = newObjectID()
	obj.waitStopped = make(chan bool, 1)

	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	queue := cellnet.NewEventQueue()
	// 创建一个tcp的侦听器，名称为server，连接地址为127.0.0.1:8801，所有连接将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	peerIns := peer.NewGenericPeer("tcp.Acceptor", name, address, queue)
	proc.BindProcessorHandler(peerIns, "tcp.ltv", obj.packetRecv)
	obj.queueIns = queue
	obj.peerIns = peerIns
	// 开始侦听
	obj.peerIns.Start()
	// 事件队列开始循环
	obj.queueIns.StartLoop()
	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	go obj.queueIns.Wait()
	return obj
}

func NewWebSocketServer(name string, address string, handle interface{}) *TcpServer {
	obj := new(TcpServer)
	obj.Name = name
	obj.Type = "websocket"
	obj.listenAddress = address
	obj.clientPool = make(map[int64]*TcpClient, 0)
	obj.clientHandle = handle
	obj.objectID = newObjectID()
	obj.waitStopped = make(chan bool, 1)

	// 创建一个事件处理队列，整个服务器只有这一个队列处理事件，服务器属于单线程服务器
	queue := cellnet.NewEventQueue()
	peerIns := peer.NewGenericPeer("gorillaws.Acceptor", name, address, queue)
	proc.BindProcessorHandler(peerIns, "gorillaws.ltv", obj.packetRecv)
	obj.queueIns = queue
	obj.peerIns = peerIns
	// 开始侦听
	obj.peerIns.Start()
	// 事件队列开始循环
	obj.queueIns.StartLoop()
	// 阻塞等待事件队列结束退出( 在另外的goroutine调用queue.StopLoop() )
	go obj.queueIns.Wait()
	return obj
}

func (self *TcpServer) String() string {
	return fmt.Sprintf("[Server][%s][%s]-%d ", self.listenAddress, self.Name, self.objectID)
}

func (self *TcpServer) WaitStop() {
	<-self.waitStopped
	LogInfo("Stopped : %v", self)
}

func (self *TcpServer) Disconnect() {
	self.Lock()
	defer func() {
		err := recover()
		if err != nil {
			LogError("Disconnect Error: %v : %v", self, err)
		}
		if self.waitStopped != nil {
			self.waitStopped <- true
			self.waitStopped = nil
		}
		self.Unlock()
	}()
	self.peerIns.Stop()
}

func (self *TcpServer) OnConnectSucc(ev cellnet.Event) {
	self.Lock()
	defer self.Unlock()

	address := ""
	if self.Type == "tcp" {
		address = ev.Session().Raw().(net.Conn).RemoteAddr().String()
	} else {
		address = ev.Session().Raw().(*websocket.Conn).RemoteAddr().String()
	}
	client := newTcpClient(address, ev)
	self.clientPool[client.SessionID()] = client

	LogInfo("Connected : %v", client)
	self.clientHandle.(TcpClientHandle).OnEventTrigger(client, "Connect")
}

func (self *TcpServer) OnDisconnect(ev cellnet.Event) {
	self.Lock()
	defer self.Unlock()

	sessionID := ev.Session().ID()
	client, ok := self.clientPool[sessionID]
	if ok {
		delete(self.clientPool, sessionID)
		LogInfo("Disconnected : %v", client)
		self.clientHandle.(TcpClientHandle).OnEventTrigger(client, "DisConnect")
	}
}

func (self *TcpServer) GetClient(sessionID int64) *TcpClient {
	self.Lock()
	defer self.Unlock()

	obj, ok := self.clientPool[sessionID]
	if ok {
		return obj
	}
	return nil
}

func (self *TcpServer) Broadcast(msg interface{}) {
	self.peerIns.(cellnet.SessionAccessor).VisitSession(
		func(ses cellnet.Session) bool {
			ses.Send(msg)
			return true
		})
}

func (self *TcpServer) packetRecv(ev cellnet.Event) {
	msg := ev.Message()
	switch msg.(type) {
	// 有新的连接
	case *cellnet.SessionAccepted:
		LogDebug("Server Accepted : %d", ev.Session().ID())
		self.OnConnectSucc(ev)
	// 有连接断开
	case *cellnet.SessionClosed:
		LogDebug("Session Closed : %d", ev.Session().ID())
		self.OnDisconnect(ev)
	default:
		client := self.GetClient(ev.Session().ID())
		if client == nil {
			LogError("Invalid Client: %v : %d", self, ev.Session().ID())
		} else {
			// 当服务器收到的是一个rpc消息
			if rpcevent, ok := ev.(*rpc.RecvMsgEvent); ok {
				response, err := self.clientHandle.(TcpClientHandle).OnRpcCommand(client, msg)
				if err != nil {
					LogError("RpcCommand Error : %v : %v", self, err)
				} else {
					if response != nil {
						rpcevent.Reply(response)
					}
				}
			} else {
				self.clientHandle.(TcpClientHandle).OnProtoCommand(client, msg)
			}
		}
	}
}
