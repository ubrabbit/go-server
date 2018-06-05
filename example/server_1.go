package main

import (
	"fmt"
)

import (
	. "github.com/ubrabbit/go-server/server"
)

var Address = "127.0.0.1:3832"

func main() {
	fmt.Println("start server:")

	InitServerPool()
	NewTcpServer("server", Address, true)

}
