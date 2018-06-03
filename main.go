package main

import (
	"fmt"

	_ "github.com/ubrabbit/go-server/common"
	config "github.com/ubrabbit/go-server/config"
	_ "github.com/ubrabbit/go-server/lib"
	_ "github.com/ubrabbit/go-server/socket"
)

func main() {
	fmt.Println("start server")

	config.InitConfig("")
}
