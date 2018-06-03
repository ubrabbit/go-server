package config

import (
	. "github.com/ubrabbit/go-server/common"
)

var g_MongoConfig *MongoConfig = nil

type MongoConfig struct {
	IP   string
	Port int
}

func GetMongoConfig() *MongoConfig {
	return g_MongoConfig
}

func InitConfigMongoDB() {
	setting := GetConfigFile().ReadConfig("mongo")
	g_MongoConfig = new(MongoConfig)
	g_MongoConfig.IP = getSettingValue(setting, "ip", 1)
	g_MongoConfig.Port = StringToInt(getSettingValue(setting, "port", 1))
}
