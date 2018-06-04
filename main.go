package main

import (
	. "github.com/ubrabbit/go-server/common"
	config "github.com/ubrabbit/go-server/config"
	_ "github.com/ubrabbit/go-server/lib"
	_ "github.com/ubrabbit/go-server/socket"
)

func main() {
	LogInfo("start server")

	config.InitConfig("settings.conf")

	LogInfo("finish server")
}
