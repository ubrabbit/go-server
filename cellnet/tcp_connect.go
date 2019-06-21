package cellnet

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
	. "github.com/ubrabbit/go-public/common"
)

const (
	RpcTimeout = 60
)

type TcpConnectHandle interface {
	OnProtoCommand(*TcpConnect, interface{})
	OnEventTrigger(*TcpConnect, string, ...interface{})
}

type TcpConnect struct {
	sync.Mutex
	Name     string
	Address  string
	queueIns cellnet.EventQueue
	peerIns  cellnet.GenericPeer

	sessionID     int64
	objectID      int64
	waitConnected chan bool

	connectHandle interface{}
}

func NewTcpConnect(name string, address string, handle interface{}) *TcpConnect {
	obj := new(TcpConnect)
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
	proc.BindProcessorHandler(peerIns, "tcp.ltv", obj.packetRecv)
	// 在peerIns接口中查询TCPConnector接口，设置连接超时1秒后自动重连
	peerIns.(cellnet.TCPConnector).SetReconnectDuration(1 * time.Second)
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
func (self *TcpConnect) String() string {
	return fmt.Sprintf("[TcpConnect][%s]-%d-%d ", self.Address, self.objectID, self.sessionID)
}

func (self *TcpConnect) ObjectID() int64 {
	return self.objectID
}

func (self *TcpConnect) Session() cellnet.Session {
	return self.peerIns.(interface {
		Session() cellnet.Session
	}).Session()
}

//sessionID在断线后通过Session获取不到
func (self *TcpConnect) SessionID() int64 {
	return self.sessionID
}

func (self *TcpConnect) Disconnect() {
	self.Lock()
	defer func() {
		err := recover()
		if err != nil {
			LogError("Disconnect Error: %v : %v", self, err)
		}
		self.Unlock()
	}()
	self.peerIns.Stop()
}

func (self *TcpConnect) OnConnectSucc(ev cellnet.Event) {
	LogInfo("ConnectSucc : %v", self)
	self.sessionID = self.Session().ID()

	//连接成功，取消阻塞
	if self.waitConnected != nil {
		self.waitConnected <- true
		self.waitConnected = nil
	}
	self.connectHandle.(TcpConnectHandle).OnEventTrigger(self, "Connect")
}

func (self *TcpConnect) OnDisconnect(ev cellnet.Event) {
	LogInfo("Disconnected : %v", self)
	self.connectHandle.(TcpConnectHandle).OnEventTrigger(self, "DisConnect")
}

func (self *TcpConnect) OnRpcAck(ev cellnet.Event) {
}

func (self *TcpConnect) PacketSend(msg interface{}) {
	self.Lock()
	defer func() {
		err := recover()
		if err != nil {
			LogError("PacketSend Error: %v : %v", self, err)
		}
		self.Unlock()
	}()
	session := self.Session()
	if session == nil {
		LogError("Session Closed: %v", self)
		return
	}
	session.Send(msg)
}

func (self *TcpConnect) RpcCall(msg interface{}) error {
	defer func() {
		err := recover()
		if err != nil {
			LogError("RpcCall Error: %v : %v", self, err)
		}
	}()
	//异步
	LogDebug("RpcCall : %v", self)
	rpc.Call(self.peerIns, msg, time.Duration(RpcTimeout)*time.Second,
		func(raw interface{}) {
			switch result := raw.(type) {
			case error:
				LogError("RpcCall Error: %v %v", self, result)
			}
		})
	return nil
}

func (self *TcpConnect) RpcCallSync(msg interface{}, callback func(*TcpConnect, interface{}, error)) error {
	defer func() {
		err := recover()
		if err != nil {
			LogError("RpcCallSync Error: %v : %v", self, err)
		}
	}()
	//同步
	LogDebug("RpcCallSync : %v", self)
	ret, err := rpc.CallSync(self.peerIns, msg, time.Duration(RpcTimeout)*time.Second)
	callback(self, ret, err)
	return err
}

func (self *TcpConnect) packetRecv(ev cellnet.Event) {
	defer func() {
		err := recover()
		if err != nil {
			LogError("packetRecv Error: %v : %v", self, err)
		}
	}()

	LogDebug("packetRecv : %v", ev.Message())
	msg := ev.Message()
	switch msg.(type) {
	case *cellnet.SessionConnected:
		self.OnConnectSucc(ev)
	case *cellnet.SessionClosed:
		self.OnDisconnect(ev)
	case *rpc.RemoteCallACK:
		self.OnRpcAck(ev)
	default:
		self.connectHandle.(TcpConnectHandle).OnProtoCommand(self, msg)
	}
}
