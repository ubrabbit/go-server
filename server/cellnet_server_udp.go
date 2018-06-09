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

type ClientHandle interface {
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
