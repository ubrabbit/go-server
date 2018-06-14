package tests

import (
	"fmt"
	"testing"
	"time"
)

import (
	. "github.com/ubrabbit/go-server/cellnet"
	proto "github.com/ubrabbit/go-server/proto"
)

var UdpAddress = "127.0.0.1:3834"

type UdpClientCmd struct {
}

func (self *UdpClientCmd) OnProtoCommand(c *UdpClient, msg interface{}) {
	switch msg := msg.(type) {
	case *proto.TestChatREQ:
		fmt.Println("TestChatREQ:  ", msg.Content)
		c.PacketSend(&proto.TestChatREQ{Content: "Udp Respond"})
	default:
		fmt.Println("Invalid Command:  ", msg)
	}
}

type UdpConnectCmd struct {
}

func (self *UdpConnectCmd) OnProtoCommand(c *UdpConnect, msg interface{}) {
	switch msg := msg.(type) {
	case *proto.TestChatREQ:
		fmt.Println("TestChatREQ:  ", msg.Content)
	default:
		fmt.Println("Invalid Command:  ", msg)
	}
}

func TestUdpServer(t *testing.T) {
	fmt.Printf("\n\n=====================  TestUdpServer  =====================\n")

	handle := &UdpClientCmd{}
	obj := NewUdpServer("server", UdpAddress, handle)
	go obj.WaitStop()
	time.Sleep(1 * time.Second)
}

func TestUdpConnect(t *testing.T) {
	fmt.Printf("\n\n=====================  TestUdpConnect  =====================\n")

	handle := &UdpConnectCmd{}
	obj := NewUdpConnect("client", UdpAddress, handle)
	time.Sleep(1 * time.Second)
	obj.PacketSend(&proto.TestChatREQ{
		Content: "TestUdpConnect",
	})
}
