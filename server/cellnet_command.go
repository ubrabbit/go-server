package server

import (
	"fmt"
)

import (
	. "github.com/ubrabbit/go-server/common"
	proto "github.com/ubrabbit/go-server/proto"
)

func onConnectCommand(c *ConnectUnit, msg interface{}) {
	LogInfo("onConnectCommand:  ", c.ObjectID())

	switch msg := msg.(type) {
	case *proto.TestChatACK:
		LogInfo("TestChatACK")
		fmt.Println("Content:  ", msg.Content)
		fmt.Println("Id:  ", msg.Id)
	case *proto.C2SConnect:
		LogInfo("C2SConnect")
		fmt.Println("msg:  ", msg)
	default:
		LogError("Invalid Command:  ", msg)
	}
}

func onServerCommand(c *ClientUnit, msg interface{}) {
	LogInfo("onServerCommand:  ", c.ObjectID())

	switch msg := msg.(type) {
	case *proto.TestChatREQ:
		LogInfo("TestChatREQ")
		msg2 := msg.Content
		fmt.Println("msg:  ", msg2)
		c.Parent.Broadcast(&proto.TestChatACK{Content: "respond222", Id: c.SessionID()})
	case *proto.C2SConnect:
		LogInfo("C2SConnect")
		fmt.Println("msg:  ", msg)
	default:
		LogError("Invalid Command:  ", msg)
	}
}
