package main

import (
	"flag"
	"fmt"
	"log/slog"
	"oauth-server/config"
	logger "oauth-server/package/log"
	"oauth-server/package/queue"
	"oauth-server/utils"
)

func init() {
	rootDir := utils.RootDir()
	configFile := fmt.Sprintf("%s/config.yml", rootDir)

	config.Init(configFile)
	logger.Init()
}

func main() {
	rabbitmqType := flag.String("type", "client", "client or server")
	flag.Parse()

	if *rabbitmqType != "client" && *rabbitmqType != "server" {
		logger.GetLogger().Info("Invalid RabbitMQ type", slog.String("type", *rabbitmqType))
		return
	}

	if *rabbitmqType == "client" {
		queue.SendDirectRabbitMQ(queue.RabbitMQQueue{
			QueueName: "HelloWorldQueue",
			Exchange:  "HelloWorldExchange",
		}, "Hello World 1")
	} else if *rabbitmqType == "server" {
		queue.ReceiveDirectRabbitMQ(queue.RabbitMQQueue{
			QueueName: "HelloWorldQueue",
			Exchange:  "HelloWorldExchange",
		})
	}
}
