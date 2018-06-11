package main

import (
	"fmt"
)

import (
	. "github.com/ubrabbit/go-common/common"
	. "github.com/ubrabbit/go-server/cellnet"
	proto "github.com/ubrabbit/go-server/proto"
)

var Address = "127.0.0.1:3832"

type ClientCmd struct {
}

func (self *ClientCmd) OnProtoCommand(c *UdpClient, msg interface{}) {
	switch msg := msg.(type) {
	case *proto.TestChatREQ:
		msg2 := msg.Content
		LogInfo("TestChatREQ:  ", msg2)
		c.PacketSend(&proto.TestChatREQ{Content: "Udp Respond"})
	case *proto.TestConnect:
		LogInfo("TestConnect:  ", msg)
	default:
		LogError("Invalid Command:  ", msg)
	}
}

func main() {
	fmt.Println("start server:")

	handle := &ClientCmd{}
	obj := NewUdpServer("server", Address, handle)
	obj.WaitStop()
}
