package server

import (
	. "github.com/ubrabbit/go-server/common"
	proto "github.com/ubrabbit/go-server/proto"
)

func onConnectCommand(c *ConnectUnit, msg interface{}) {
	LogInfo("onConnectCommand:  ", c.ObjectID())

	switch msg := msg.(type) {
	case *proto.TestChatACK:
		LogInfo("TestChatACK:  ", msg)
	case *proto.C2SConnect:
		LogInfo("C2SConnect:  ", msg)
	default:
		LogError("Invalid Command:  ", msg)
	}
}

func onServerCommand(c *ClientUnit, msg interface{}) {
	//LogInfo("onServerCommand:  ", c.ObjectID())
	switch msg := msg.(type) {
	case *proto.TestChatREQ:
		msg2 := msg.Content
		LogInfo("TestChatREQ:  ", msg2)
		c.Parent.Broadcast(&proto.TestChatACK{Content: "respond:  " + msg2, Id: c.SessionID()})
	case *proto.C2SConnect:
		LogInfo("C2SConnect:  ", msg)
	default:
		LogError("Invalid Command:  ", msg)
	}
}
