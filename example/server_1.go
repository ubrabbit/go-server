package main

import (
	"fmt"
	"time"
)

import (
	. "github.com/ubrabbit/go-server/common"
	proto "github.com/ubrabbit/go-server/proto"
	. "github.com/ubrabbit/go-server/server"
)

var Address = "127.0.0.1:3832"

type ClientCmd struct {
}

func (self *ClientCmd) OnProtoCommand(c *TcpClient, msg interface{}) {
	//LogInfo("onServerCommand:  ", c.ObjectID())
	switch msg := msg.(type) {
	case *proto.TestChatREQ:
		msg2 := msg.Content
		LogInfo("TestChatREQ:  ", msg2)
		c.Parent.Broadcast(&proto.TestChatACK{Content: "respond_start: " + msg2 + " finish_response", Id: c.SessionID()})
	case *proto.C2SConnect:
		LogInfo("C2SConnect:  ", msg)
	default:
		LogError("Invalid Command:  ", msg)
	}
}

func (self *ClientCmd) OnRpcCommand(c *TcpClient, msg interface{}) (interface{}, error) {
	fmt.Println(">>>>>>>>>>>>>>>>>>>> OnRpcCommand")
	time.Sleep(10 * time.Second)
	fmt.Println(">>>>>>>>>>>>>>>>>>>> OnRpcCommand Ack")
	return proto.TestChatREQ{Content: "Rpc_Respond"}, nil
}

func (self *ClientCmd) OnEventTrigger(c *TcpClient, name string, args ...interface{}) {
	LogInfo(c, " EventTrigger: ", name, args)
}

func main() {
	fmt.Println("start server:")

	handle := &ClientCmd{}
	obj := NewTcpServer("server", Address, handle)
	obj.WaitStop()
}
