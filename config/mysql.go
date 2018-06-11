package config

var g_MysqlConfig *MysqlConfig = nil

type MysqlConfig struct {
	IP       string
	Port     int
	Account  string
	Password string
	Database string
}

func GetMysqlConfig() *MysqlConfig {
	return g_MysqlConfig
}

func InitConfigMysql() {
	setting := GetConfigFile().ReadConfig("mysql")

	g_MysqlConfig = new(MysqlConfig)
	g_MysqlConfig.IP = getSettingValue(setting, "ip", 1)
	g_MysqlConfig.Port = stringToInt(getSettingValue(setting, "port", 1))
	g_MysqlConfig.Account = getSettingValue(setting, "account", 1)
	g_MysqlConfig.Password = getSettingValue(setting, "password", 1)
	g_MysqlConfig.Database = getSettingValue(setting, "database", 1)
}
