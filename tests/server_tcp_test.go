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
}

func TestTcpConnect(t *testing.T) {
	fmt.Printf("\n\n=====================  TestTcpConnect  =====================\n")

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
