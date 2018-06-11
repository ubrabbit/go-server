package config

var g_RedisConfig *RedisConfig = nil

type RedisConfig struct {
	IP   string
	Port int
}

func GetRedisConfig() *RedisConfig {
	return g_RedisConfig
}

func InitConfigRedis() {
	setting := GetConfigFile().ReadConfig("redis")

	g_RedisConfig = new(RedisConfig)
	g_RedisConfig.IP = getSettingValue(setting, "ip", 1)
	g_RedisConfig.Port = stringToInt(getSettingValue(setting, "port", 1))
}
