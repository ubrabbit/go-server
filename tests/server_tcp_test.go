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

var Address = "127.0.0.1:3832"

func TestTcpServer(t *testing.T) {
	fmt.Printf("\n\n=====================  TestTcpServer  =====================\n")

	handle := &ClientCmd{}
	obj := NewTcpServer("server", Address, handle)
	go obj.WaitStop()

	go func() {
		time.Sleep(2 * time.Second)
		obj.Disconnect()
	}()
}

func TestTcpConnect(t *testing.T) {
	fmt.Printf("\n\n=====================  TestTcpConnect  =====================\n")

	handle := &ConnectCmd{}
	obj := NewTcpConnect("client", Address, handle)
	msg := proto.TestConnect{
		Hello:    "TestConnect 1",
		Account:  "TestConnect 2",
		Password: "TestConnect 3",
	}

	time.Sleep(1 * time.Second)
	obj.PacketSend(&msg)
	time.Sleep(2 * time.Second)
	obj.PacketSend(&msg)
}
