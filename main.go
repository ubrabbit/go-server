package main

import (
	. "github.com/ubrabbit/go-server/common"
	config "github.com/ubrabbit/go-server/config"
	server "github.com/ubrabbit/go-server/config"
)

func main() {
	LogInfo("start server")

	config.InitConfig("settings.conf")
	server.InitServer()

	LogInfo("finish server")
}
