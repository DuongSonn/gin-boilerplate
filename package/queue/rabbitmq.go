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

var (
	RABBIT_MQ_QUEUE_REGISTER = "register"
	CLIENT_EXCHANGE          = map[string]rabbitMQExchange{
		"RPCClientRoutingKey": {
			Exchange:   "RPCClientExchange",
			RoutingKey: "RPCClientRoutingKey",
		},
	}
)

type RabbitMQQueue struct {
	QueueName  string
	Exchange   string
	RoutingKey string
	Consumer   string
}

type RabbitMQRPCQueue struct {
	Server RabbitMQQueue
	Client RabbitMQQueue
}

type rabbitMQConnect struct {
	ch   *amqp.Channel
	conn *amqp.Connection
}
type rabbitMQExchange struct {
	Exchange   string
	RoutingKey string
}

func connectRabbitMQ() (*rabbitMQConnect, error) {
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

	return &rabbitMQConnect{ch, conn}, nil
}

func getExchangeName(routingKey string) string {
	if exchange, ok := CLIENT_EXCHANGE[routingKey]; ok {
		return exchange.Exchange
	}

	return ""
}

func SendRPCRabbitMQ(queueConf RabbitMQRPCQueue) {
	connect, err := connectRabbitMQ()
	if err != nil {
		logger.GetLogger().Info("Connect RabbitMQ", slog.String("error", err.Error()))
		return
	}
	defer func() {
		connect.conn.Close()
		connect.ch.Close()
	}()
	ch := connect.ch

	// Declare queue
	if _, err := ch.QueueDeclare(
		queueConf.Client.QueueName, // name
		false,                      // durable
		false,                      // delete when unused
		true,                       // exclusive
		false,                      // noWait
		nil,                        // arguments
	); err != nil {
		logger.GetLogger().Info("QueueDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare server exchange
	if err := ch.ExchangeDeclare(
		queueConf.Server.Exchange, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	); err != nil {
		logger.GetLogger().Info("ExchangeDeclare Server", slog.String("error", err.Error()))
		return
	}
	if err := ch.ExchangeDeclare(
		queueConf.Client.Exchange, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	); err != nil {
		logger.GetLogger().Info("ExchangeDeclare Client", slog.String("error", err.Error()))
		return
	}
	// Bind client queue to client exchange
	if err = ch.QueueBind(
		queueConf.Client.QueueName,  // queue name
		queueConf.Client.RoutingKey, // routing key
		queueConf.Client.Exchange,   // exchange
		false,
		nil,
	); err != nil {
		logger.GetLogger().Info("QueueBind RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare consumer
	msgs, err := ch.Consume(
		queueConf.Client.QueueName, // queue
		queueConf.Client.Consumer,  // consumer
		true,                       // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
	)
	if err != nil {
		logger.GetLogger().Info("Consume RabbitMQ", slog.String("error", err.Error()))
		return
	}

	ID := utils.GenerateUUID().String()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = ch.PublishWithContext(ctx,
		queueConf.Server.Exchange,   // exchange
		queueConf.Server.RoutingKey, // routing key
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: ID,
			ReplyTo:       queueConf.Client.RoutingKey,
			Body:          []byte("Test RPC"),
			Expiration:    "5000", // 5 seconds
		},
	); err != nil {
		logger.GetLogger().Info("PublishWithContext RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Wait for response
	select {
	case d := <-msgs:
		if ID == d.CorrelationId {
			fmt.Println(" [x] Received ", string(d.Body))
			return // Exit on successful response
		}
	case <-ctx.Done():
		logger.GetLogger().Info("Timeout waiting for RPC response")
		return // Exit on timeout
	}
}

func ReceiveRPCRabbitMQ(queueConf RabbitMQRPCQueue) {
	connect, err := connectRabbitMQ()
	if err != nil {
		logger.GetLogger().Info("Connect RabbitMQ", slog.String("error", err.Error()))
		return
	}
	defer func() {
		connect.conn.Close()
		connect.ch.Close()
	}()
	ch := connect.ch

	// Declare queue
	if _, err := ch.QueueDeclare(
		queueConf.Server.QueueName, // name
		false,                      // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	); err != nil {
		logger.GetLogger().Info("QueueDeclare RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare server exchange
	if err = ch.ExchangeDeclare(
		queueConf.Server.Exchange, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	); err != nil {
		logger.GetLogger().Info("ExchangeDeclare Sever", slog.String("error", err.Error()))
		return
	}
	// Declare client exchange
	for _, exchange := range CLIENT_EXCHANGE {
		if err = ch.ExchangeDeclare(
			exchange.Exchange, // name
			"direct",          // type
			true,              // durable
			false,             // auto-deleted
			false,             // internal
			false,             // no-wait
			nil,               // arguments
		); err != nil {
			logger.GetLogger().Info("ExchangeDeclare Client", slog.String("error", err.Error()))
			return
		}
	}

	// Bind server queue to server exchange
	if err = ch.QueueBind(
		queueConf.Server.QueueName,  // queue name
		queueConf.Server.RoutingKey, // routing key
		queueConf.Server.Exchange,   // exchange
		false,
		nil,
	); err != nil {
		logger.GetLogger().Info("QueueBind RabbitMQ", slog.String("error", err.Error()))
		return
	}

	// Declare consumer
	msgs, err := ch.Consume(
		queueConf.Server.QueueName, // queue
		queueConf.Server.Consumer,  // consumer
		false,                      // auto-ack
		false,                      // exclusive
		false,                      // no-local
		false,                      // no-wait
		nil,                        // args
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
			routingKey := d.ReplyTo
			exchange := getExchangeName(routingKey)

			if err := ch.PublishWithContext(ctx,
				exchange,   // exchange
				routingKey, // routing key
				false,      // mandatory
				false,      // immediate
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
	connect, err := connectRabbitMQ()
	if err != nil {
		logger.GetLogger().Info("Connect RabbitMQ", slog.String("error", err.Error()))
		return
	}
	defer func() {
		connect.conn.Close()
		connect.ch.Close()
	}()
	ch := connect.ch

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
	connect, err := connectRabbitMQ()
	if err != nil {
		logger.GetLogger().Info("Connect RabbitMQ", slog.String("error", err.Error()))
		return
	}
	defer func() {
		connect.conn.Close()
		connect.ch.Close()
	}()
	ch := connect.ch

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
