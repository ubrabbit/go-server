package tests

import (
	"fmt"
	. "github.com/ubrabbit/go-server/server"
	"testing"
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
	NewTcpClient("client", Address)
}
