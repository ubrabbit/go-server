package config

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
	g_MongoConfig.Port = stringToInt(getSettingValue(setting, "port", 1))
}
