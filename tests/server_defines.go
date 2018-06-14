package tests

import (
	"fmt"
)

import (
	. "github.com/ubrabbit/go-server/cellnet"
	proto "github.com/ubrabbit/go-server/proto"
)

type ClientCmd struct {
}

func (self *ClientCmd) OnProtoCommand(c *TcpClient, msg interface{}) {
	switch msg := msg.(type) {
	case *proto.TestConnect:
		fmt.Println("Server Recv TestConnect:  ", msg)
		c.PacketSend(&proto.TestChatREQ{Content: "TestConnect Respond"})
	case *proto.TestWSJson:
		c.PacketSend(&proto.TestWSJson{Msg: "TestWSJson Respond", Value: 10086})
		fmt.Println("Server Recv TestWSJson:  ", msg)
	default:
		fmt.Println("Invalid Command:  ", msg)
	}
}

func (self *ClientCmd) OnRpcCommand(c *TcpClient, msg interface{}) (interface{}, error) {
	return proto.TestChatREQ{Content: "OnRpcCommand"}, nil
}

func (self *ClientCmd) OnEventTrigger(c *TcpClient, name string, args ...interface{}) {
}

type ConnectCmd struct {
}

func (self *ConnectCmd) OnProtoCommand(c *TcpConnect, msg interface{}) {
}

func (self *ConnectCmd) OnEventTrigger(c *TcpConnect, name string, args ...interface{}) {
}
