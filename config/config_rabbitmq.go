package config

var g_RabbitMQConfig *RabbitMQConfig = nil

type RabbitMQConfig struct {
	Account  string
	Password string
	Host     string
	HostName string
	Port     int
}

func GetRabbitMQConfig() *RabbitMQConfig {
	return g_RabbitMQConfig
}

func InitConfigRabbitMQ() {
	setting := GetConfigFile().ReadConfig("rabbitmq")

	g_RabbitMQConfig = new(RabbitMQConfig)
	g_RabbitMQConfig.Account = getSettingValue(setting, "account", 1)
	g_RabbitMQConfig.Password = getSettingValue(setting, "password", 1)
	g_RabbitMQConfig.Host = getSettingValue(setting, "host", 1)
	g_RabbitMQConfig.HostName = getSettingValue(setting, "hostname", 1)
	g_RabbitMQConfig.Port = stringToInt(getSettingValue(setting, "port", 1))
}
