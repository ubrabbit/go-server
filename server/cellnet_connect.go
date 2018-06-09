package server

import (
	"fmt"
	"sync"
	"time"
)

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/peer"
	_ "github.com/davyxu/cellnet/peer/tcp"
	"github.com/davyxu/cellnet/proc"
	_ "github.com/davyxu/cellnet/proc/tcp"
	"github.com/davyxu/cellnet/rpc"
)

import (
	. "github.com/ubrabbit/go-server/common"
)

const (
	RpcTimeout = 60
)

type ConnectHandle interface {
	OnProtoCommand(*Connect, interface{})
	OnEventTrigger(*Connect, string, ...interface{})
}

type Connect struct {
	sync.Mutex
	Name    string
	Address string
	Queue   cellnet.EventQueue
	Peer    cellnet.GenericPeer

	sessionID     int64
	objectID      int64
	waitConnected chan bool

	connectHandle interface{}
}

func NewTcpConnect(name string, address string, handle interface{}) *Connect {
	obj := new(Connect)
	obj.Name = name
	obj.Address = address
	obj.objectID = newObjectID()
	obj.sessionID = 0
	obj.connectHandle = handle
	obj.waitConnected = make(chan bool, 1)

	// 创建一个事件处理队列，整个客户端只有这一个队列处理事件，客户端属于单线程模型
	queue := cellnet.NewEventQueue()
	// 创建一个tcp的连接器，名称为Connect，连接地址为127.0.0.1:8801，将事件投递到queue队列,单线程的处理（收发封包过程是多线程）
	// peer.NewGenericPeer("tcp.Connector", "Connect", "127.0.0.1:18801", queue)
	peerIns := peer.NewGenericPeer("tcp.Connector", name, address, queue)
	proc.BindProcessorHandler(peerIns, "tcp.ltv", obj.PacketRecv)
	// 在peerIns接口中查询TCPConnector接口，设置连接超时1秒后自动重连
	peerIns.(cellnet.TCPConnector).SetReconnectDuration(1 * time.Second)
	obj.Queue = queue
	obj.Peer = peerIns
	// 开始发起到服务器的连接
	obj.Peer.Start()
	// 事件队列开始循环
	obj.Queue.StartLoop()

	// 等待连接成功再返回
	<-obj.waitConnected
	obj.waitConnected = nil
	return obj
}

//__repr__
func (self *Connect) String() string {
	return fmt.Sprintf("[Connect][%s]-%d-%d ", self.Address, self.objectID, self.sessionID)
}

func (self *Connect) ObjectID() int64 {
	return self.objectID
}

func (self *Connect) Session() cellnet.Session {
	return self.Peer.(interface {
		Session() cellnet.Session
	}).Session()
}

//sessionID在断线后通过Session获取不到
func (self *Connect) SessionID() int64 {
	return self.sessionID
}

func (self *Connect) Disconnect() {
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

func (self *Connect) OnConnectSucc(ev cellnet.Event) {
	self.Lock()
	defer self.Unlock()

	LogInfo(self, "ConnectSucc")
	self.sessionID = self.Session().ID()
	pool := GetConnectPool()
	pool.Add(self)

	//连接成功，取消阻塞
	self.waitConnected <- true
	self.connectHandle.(ConnectHandle).OnEventTrigger(self, "Connect")
}

func (self *Connect) OnDisconnect(ev cellnet.Event) {
	self.Lock()
	defer self.Unlock()

	LogInfo(self, "Disconnected")
	pool := GetConnectPool()
	pool.Remove(self.SessionID())

	self.connectHandle.(ConnectHandle).OnEventTrigger(self, "DisConnect")
}

func (self *Connect) OnRpcAck(ev cellnet.Event) {
}

func (self *Connect) PacketSend(msg interface{}) {
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

func (self *Connect) PacketRecv(ev cellnet.Event) {
	defer func() {
		err := recover()
		if err != nil {
			LogError(self, "PacketRecv: ", err)
		}
	}()

	LogInfo("PacketRecv")
	msg := ev.Message()
	switch msg.(type) {
	case *cellnet.SessionConnected:
		self.OnConnectSucc(ev)
	case *cellnet.SessionClosed:
		self.OnDisconnect(ev)
	case *rpc.RemoteCallACK:
		self.OnRpcAck(ev)
	default:
		self.connectHandle.(ConnectHandle).OnProtoCommand(self, msg)
	}
}

func (self *Connect) RpcCall(msg interface{}) error {
	defer func() {
		err := recover()
		if err != nil {
			LogError(self, "RpcCall Error: ", err)
		}
	}()
	//异步
	LogInfo("RpcCall")
	rpc.Call(self.Peer, msg, time.Duration(RpcTimeout)*time.Second,
		func(raw interface{}) {
			switch result := raw.(type) {
			case error:
				LogError(self, "RpcCall Error: ", result)
			}
		})
	return nil
}

func (self *Connect) RpcCallSync(msg interface{}, callback func(*Connect, interface{}, error)) error {
	defer func() {
		err := recover()
		if err != nil {
			LogError(self, "RpcCallSync Error: ", err)
		}
	}()
	//同步
	LogInfo("RpcCallSync")
	ret, err := rpc.CallSync(self.Peer, msg, time.Duration(RpcTimeout)*time.Second)
	callback(self, ret, err)
	return err
}
