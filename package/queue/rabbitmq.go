package queue

import (
	"context"
	"fmt"
	"log"
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

type RabbitMQQueue struct {
	QueueName  string
	Exchange   string
	RoutingKey string
	Consumer   string
}

func connectRabbitMQ() (*amqp.Channel, error) {
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

func SendRPCRabbitMQ(queueConf RabbitMQQueue) {
	ch, err := connectRabbitMQ()
	if err != nil {
		logger.GetLogger().Info("Connect RabbitMQ", slog.String("error", err.Error()))
		return
	}
	defer ch.Close()

	// Declare queue
	if _, err := ch.QueueDeclare(
		queueConf.QueueName, // name
		false,               // durable
		false,               // delete when unused
		true,                // exclusive
		false,               // noWait
		nil,                 // arguments
	); err != nil {
		logger.GetLogger().Info("QueueDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare exchange
	if err := ch.ExchangeDeclare(
		queueConf.Exchange, // name
		"fanout",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	); err != nil {
		logger.GetLogger().Info("ExchangeDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare consumer
	msgs, err := ch.Consume(
		queueConf.QueueName, // queue
		queueConf.Consumer,  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
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

	if err = ch.PublishWithContext(ctx,
		queueConf.Exchange, // exchange
		"",                 // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: ID,
			ReplyTo:       queueConf.QueueName,
			Body:          []byte("Test RPC"),
		},
	); err != nil {
		logger.GetLogger().Info("PublishWithContext RabbitMQ", slog.String("error", err.Error()))
		return
	}

	for d := range msgs {
		if ID == d.CorrelationId {
			fmt.Println(" [x] Received ", string(d.Body))
			d.Ack(false)
			break
		}
	}
}

func ReceiveRPCRabbitMQ(queueConf RabbitMQQueue) {
	ch, err := connectRabbitMQ()
	if err != nil {
		logger.GetLogger().Info("Connect RabbitMQ", slog.String("error", err.Error()))
		return
	}
	defer ch.Close()

	// Declare queue
	if _, err := ch.QueueDeclare(
		queueConf.QueueName, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	); err != nil {
		logger.GetLogger().Info("QueueDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare exchange
	if err = ch.ExchangeDeclare(
		queueConf.Exchange, // name
		"fanout",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	); err != nil {
		logger.GetLogger().Info("ExchangeDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Bind queue to exchange
	if err = ch.QueueBind(
		queueConf.QueueName,  // queue name
		queueConf.RoutingKey, // routing key
		queueConf.Exchange,   // exchange
		false,
		nil,
	); err != nil {
		logger.GetLogger().Info("QueueBind RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare consumer
	msgs, err := ch.Consume(
		queueConf.QueueName, // queue
		queueConf.Consumer,  // consumer
		false,               // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		logger.GetLogger().Info("Consume RabbitMQ", slog.String("error", err.Error()))
		return
	}

	var forever chan struct{}

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)

			if err := ch.PublishWithContext(ctx,
				queueConf.Exchange, // exchange
				d.ReplyTo,          // routing key
				false,              // mandatory
				false,              // immediate
				amqp.Publishing{
					ContentType:   "text/plain",
					CorrelationId: d.CorrelationId,
					Body:          []byte("Response RPC Test"),
				},
			); err != nil {
				logger.GetLogger().Info("PublishWithContext RabbitMQ", slog.String("error", err.Error()))
			}

			d.Ack(false)
		}
	}()

	log.Printf(" [*] Awaiting RPC requests")
	<-forever
}

func SendDirectRabbitMQ(queueConf RabbitMQQueue, message string) {
	ch, err := connectRabbitMQ()
	if err != nil {
		logger.GetLogger().Info("Connect RabbitMQ", slog.String("error", err.Error()))
		return
	}
	defer ch.Close()

	// Declare queue
	if _, err := ch.QueueDeclare(
		queueConf.QueueName, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	); err != nil {
		logger.GetLogger().Info("QueueDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare exchange
	if err := ch.ExchangeDeclare(
		queueConf.Exchange, // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	); err != nil {
		logger.GetLogger().Info("ExchangeDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ch.PublishWithContext(
		ctx,
		queueConf.Exchange,   // exchange
		queueConf.RoutingKey, // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	); err != nil {
		logger.GetLogger().Info("PublishWithContext RabbitMQ", slog.String("error", err.Error()))
		return
	}

	log.Printf(" [x] Sent %s\n", message)
}

func ReceiveDirectRabbitMQ(queueConf RabbitMQQueue) {
	ch, err := connectRabbitMQ()
	if err != nil {
		logger.GetLogger().Info("Connect RabbitMQ", slog.String("error", err.Error()))
		return
	}
	defer ch.Close()

	// Declare queue

	if _, err := ch.QueueDeclare(
		queueConf.QueueName, // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	); err != nil {
		logger.GetLogger().Info("QueueDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare exchange
	if err = ch.ExchangeDeclare(
		queueConf.Exchange, // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	); err != nil {
		logger.GetLogger().Info("ExchangeDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Bind queue to exchange

	if err = ch.QueueBind(
		queueConf.QueueName,  // queue name
		queueConf.RoutingKey, // routing key
		queueConf.Exchange,   // exchange
		false,
		nil,
	); err != nil {
		logger.GetLogger().Info("QueueBind RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare consumer
	msgs, err := ch.Consume(
		queueConf.QueueName, // queue
		"",                  // consumer
		true,                // auto-ack
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		logger.GetLogger().Info("Consume RabbitMQ", slog.String("error", err.Error()))
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()
	<-forever
}
