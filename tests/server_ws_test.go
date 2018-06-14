package tests

import (
	"fmt"
	"testing"
	"time"
)

import (
	. "github.com/ubrabbit/go-server/cellnet"
)

var WSAddress = "http://127.0.0.1:3833/test"

func TestWSServer(t *testing.T) {
	fmt.Printf("\n\n=====================  TestWSServer  =====================\n")

	handle := &ClientCmd{}
	obj := NewWebSocketServer("server", WSAddress, handle)
	go obj.WaitStop()
	time.Sleep(10 * time.Second)
}
