package tests

import (
	config "github.com/ubrabbit/go-server/config"
	"testing"
)

const filepath = "config_test.conf"

func TestConfig(t *testing.T) {
	config.InitConfig(filepath)
	config.InitConfigMongoDB()
	config.InitConfigMysql()
	config.InitConfigRabbitMQ()
	config.InitConfigRedis()
}
