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
	serverType := flag.String("server", "normal", "normal or rpc")
	flag.Parse()

	if *rabbitmqType != "client" && *rabbitmqType != "server" {
		logger.GetLogger().Info("Invalid RabbitMQ type", slog.String("type", *rabbitmqType))
		return
	}

	if *serverType != "normal" && *serverType != "rpc" {
		logger.GetLogger().Info("Invalid server type", slog.String("type", *serverType))
		return
	}

	switch *serverType {
	case "normal":
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
	case "rpc":
		if *rabbitmqType == "client" {
			queue.SendRPCRabbitMQ(queue.RabbitMQRPCQueue{
				Client: queue.RabbitMQQueue{
					QueueName:  "RPCClientQueue",
					Exchange:   "RPCClientExchange",
					RoutingKey: "RPCClientRoutingKey",
					Consumer:   "RPCClientConsumer",
				},
				Server: queue.RabbitMQQueue{
					QueueName:  "RPCServerQueue",
					Exchange:   "RPCServerExchange",
					RoutingKey: "RPCServerRoutingKey",
					Consumer:   "RPCServerConsumer",
				},
			})
		} else if *rabbitmqType == "server" {
			queue.ReceiveRPCRabbitMQ(queue.RabbitMQRPCQueue{
				Client: queue.RabbitMQQueue{
					QueueName:  "RPCClientQueue",
					Exchange:   "RPCClientExchange",
					RoutingKey: "RPCClientRoutingKey",
					Consumer:   "RPCClientConsumer",
				},
				Server: queue.RabbitMQQueue{
					QueueName:  "RPCServerQueue",
					Exchange:   "RPCServerExchange",
					RoutingKey: "RPCServerRoutingKey",
					Consumer:   "RPCServerConsumer",
				},
			})
		}
	}

}
