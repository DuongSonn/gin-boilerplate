package queue

import (
	"context"
	"fmt"
	"log/slog"
	"oauth-server/config"
	logger "oauth-server/package/log"
	"oauth-server/utils"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	RABBIT_MQ_QUEUE_REGISTER = "register"
)

func connect() (*amqp.Channel, error) {
	conf := config.GetConfiguration().RabbitMQ
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", conf.User, conf.Password, conf.Host, conf.Port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func RequestRPCRabbitMQ(queueName string) {
	ch, err := connect()
	if err != nil {
		logger.GetLogger().Info("Connect RabbitMQ", slog.String("error", err.Error()))
		return
	}
	defer ch.Close()

	// Declare queue
	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // noWait
		nil,       // arguments
	)
	if err != nil {
		logger.GetLogger().Info("QueueDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare consumer
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logger.GetLogger().Info("Consume RabbitMQ", slog.String("error", err.Error()))
		return
	}

	ID := utils.GenerateUUID().String()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer func() {
		fmt.Println(" [x] Done")
		cancel()
	}()

	err = ch.PublishWithContext(ctx,
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: ID,
			ReplyTo:       q.Name,
			Body:          []byte("test"),
		})
	if err != nil {
		logger.GetLogger().Info("PublishWithContext RabbitMQ", slog.String("error", err.Error()))
		return
	}

	for d := range msgs {
		if ID == d.CorrelationId {
			fmt.Println(" [x] Received ", string(d.Body))
			break
		}
	}
}
