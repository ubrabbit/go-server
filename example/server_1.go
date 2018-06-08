package main

import (
	"fmt"
)

import (
	. "github.com/ubrabbit/go-server/common"
	proto "github.com/ubrabbit/go-server/proto"
	. "github.com/ubrabbit/go-server/server"
)

var Address = "127.0.0.1:3832"

func onServerCommand(c *ClientUnit, msg interface{}) {
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

func onEventTrigger(c *ClientUnit, name string, args ...interface{}) {
	LogInfo(c, " EventTrigger: ", name, args)
}

func main() {
	fmt.Println("start server:")

	InitServerPool()
	NewTcpServer("server", Address, true, onServerCommand, onEventTrigger)
}
