package server

import (
	"fmt"
)

import (
	"github.com/davyxu/cellnet"
	. "github.com/ubrabbit/go-server/common"
	proto "github.com/ubrabbit/go-server/proto"
)

func onClientCommand(ev cellnet.Event) {
	switch msg := ev.Message().(type) {
	case proto.C2SConnect:
		LogInfo("S2CConnect")
		fmt.Println("msg:  ", msg)
	}
}

func onServerCommand(ev cellnet.Event) {
	LogInfo("onServerCommand")

	switch msg := ev.Message().(type) {
	case *proto.C2SConnect:
		LogInfo("C2SConnect")
		fmt.Println("msg:  ", msg)
	case *proto.TestChatREQ:
		LogInfo("TestChatREQ")
		msg2 := msg.Content
		fmt.Println("msg:  ", msg2)
	}
}
