package main

import (
	"fmt"
	"time"
)

import (
	. "github.com/ubrabbit/go-common/common"
	. "github.com/ubrabbit/go-server/cellnet"
	proto "github.com/ubrabbit/go-server/proto"
)

var Address = "127.0.0.1:3832"

type ClientCmd struct {
}

func (self *ClientCmd) OnProtoCommand(c *TcpClient, msg interface{}) {
	switch msg := msg.(type) {
	case *proto.TestChatREQ:
		LogInfo("TestChatREQ:  %v", msg.Content)
	case *proto.TestConnect:
		LogInfo("TestConnect:  %v", msg)
	default:
		LogError("Invalid Command:  %v", msg)
	}
}

func (self *ClientCmd) OnRpcCommand(c *TcpClient, msg interface{}) (interface{}, error) {
	fmt.Println(">>>>>>>>>>>>>>>>>>>> OnRpcCommand")
	time.Sleep(3 * time.Second)
	return proto.TestChatREQ{Content: "Rpc_Respond"}, nil
}

func (self *ClientCmd) OnEventTrigger(c *TcpClient, name string, args ...interface{}) {
	LogInfo(" EventTrigger: %v %s %v", c, name, args)
}

func main() {
	fmt.Println("start server:")

	handle := &ClientCmd{}
	obj := NewTcpServer("server", Address, handle)
	obj.WaitStop()
}
