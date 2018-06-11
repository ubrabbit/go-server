package tests

import (
	"fmt"
	"testing"
	"time"
)

import (
	proto "github.com/ubrabbit/go-server/proto"
	. "github.com/ubrabbit/go-server/server"
)

var Address = "127.0.0.1:3832"

type ClientCmd struct {
}

func (self *ClientCmd) OnProtoCommand(c *TcpClient, msg interface{}) {
	switch msg := msg.(type) {
	case *proto.TestConnect:
		fmt.Println("Server Recv TestConnect:  ", msg)
	default:
		fmt.Println("Invalid Command:  ", msg)
	}
}

func (self *ClientCmd) OnRpcCommand(c *TcpClient, msg interface{}) (interface{}, error) {
	return proto.TestChatREQ{Content: "OnRpcCommand"}, nil
}

func (self *ClientCmd) OnEventTrigger(c *TcpClient, name string, args ...interface{}) {
	fmt.Println(c, " OnEventTrigger: ", name, args)
}

type ConnectCmd struct {
}

func (self *ConnectCmd) OnProtoCommand(c *TcpConnect, msg interface{}) {
}

func (self *ConnectCmd) OnEventTrigger(c *TcpConnect, name string, args ...interface{}) {
	fmt.Println("OnEventTrigger:  ", c, name, args)
}

func TestServer(t *testing.T) {
	fmt.Printf("\n\n=====================  TestServer  =====================\n")

	handle := &ClientCmd{}
	obj := NewTcpServer("server", Address, handle)
	go obj.WaitStop()
}

func TestClient(t *testing.T) {
	fmt.Printf("\n\n=====================  TestClient  =====================\n")

	a := "hello"
	b := "aaa"
	c := "bbb"
	handle := &ConnectCmd{}
	obj := NewTcpConnect("client", Address, handle)
	msg := proto.TestConnect{
		Hello:    a,
		Account:  b,
		Password: c,
	}
	obj.PacketSend(&msg)

	time.Sleep(1 * time.Second)
}
