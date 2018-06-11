package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	g_RunPath   string = ""
	g_LogPath   string = ""
	g_CachePath string = ""
	g_DataPath  string = ""
)

func GetPath(path string) string {
	switch path {
	case "run":
		return g_RunPath
	case "log":
		return g_LogPath
	case "data":
		return g_DataPath
	case "cache":
		return g_CachePath
	default:
		return g_CachePath
	}
}

func InitRunPath() {
	cfg := GetConfigFile()
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if cfg.HasOption("system", "run") {
		path = cfg.ReadConfigOption("system", "run")
	}
	g_RunPath = path
	fmt.Println("Run path: ", g_RunPath)

	g_LogPath = strings.Join([]string{g_RunPath, "log"}, "/")
	err = os.MkdirAll(g_LogPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	g_CachePath = strings.Join([]string{g_RunPath, "cache"}, "/")
	err = os.MkdirAll(g_CachePath, 0755)
	if err != nil {
		log.Fatal(err)
	}
	g_DataPath = strings.Join([]string{g_RunPath, "data"}, "/")
	err = os.MkdirAll(g_DataPath, 0755)
	if err != nil {
		log.Fatal(err)
	}
}
