package tests

import (
	"fmt"
	"testing"
)

import (
	proto "github.com/ubrabbit/go-server/proto"
	. "github.com/ubrabbit/go-server/server"
)

var Address = "127.0.0.1:3832"

func TestServer(t *testing.T) {
	fmt.Printf("\n\n=====================  TestServer  =====================\n")

	InitServerPool()
	NewTcpServer("server", Address, false)
}

func TestClient(t *testing.T) {
	fmt.Printf("\n\n=====================  TestClient  =====================\n")

	InitClientPool()
	a := "hello"
	b := "aaa"
	c := "bbb"
	obj := NewTcpClient("client", Address)
	msg := proto.C2SConnect{
		Hello:    &a,
		Account:  &b,
		Password: &c,
	}
	obj.PacketSend(msg)
}
