package config

import (
	"log"

	"github.com/spf13/viper"
)

var config *Configuration

type Configuration struct {
	MysqlDB    MysqlDatabase    `mapstructure:"mysql"`
	PostgresDB PostgresDatabase `mapstructure:"postgres"`
	Server     Server           `mapstructure:"server"`
	Jwt        JWT              `mapstructure:"jwt"`
	GRPC       GRPC             `mapstructure:"grpc"`
	RabbitMQ   RabbitMQ         `mapstructure:"rabbitmq"`
}

func Init(filePath string) {
	var configuration *Configuration

	viper.SetConfigFile(filePath)
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("error_decode_config: %v", err)
	}

	if configuration.Server.Port == 0 {
		configuration.Server.Port = 8080
	}

	config = configuration
}

func GetConfiguration() *Configuration {
	return config
}
