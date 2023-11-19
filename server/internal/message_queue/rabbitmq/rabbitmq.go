package rabbitmq

import (
	"context"
	"fmt"
	"logswift/internal/app/config"
	messagehandler "logswift/internal/message_queue/message_handler"
	"logswift/pkg/logger"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQClient struct {
	client *amqp091.Connection
	logger logger.ILogger
	name   string
	ch     *amqp091.Channel
	queue  amqp091.Queue
}

func NewRabbitMQClient() *RabbitMQClient {
	return &RabbitMQClient{
		logger: logger.GetInstance(),
	}
}

func (mq *RabbitMQClient) Connect(cfg config.MQConfig) error {
	mq.logger.Info("connecting to rabbitmq", "config", cfg)
	var err error
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	mq.client, err = amqp091.Dial(url)
	if err != nil {
		mq.logger.Error("error connecting to rabbitmq", "error", err.Error())
		return err
	}

	mq.logger.Info("connected to rabbitmq successfully")
	mq.name = cfg.Name
	return nil
}

func (mq *RabbitMQClient) DeclareQueue() error {
	var err error
	mq.ch, err = mq.client.Channel()
	if err != nil {
		mq.logger.Error("error creating channel", "error", err.Error())
		return err
	}

	mq.queue, err = mq.ch.QueueDeclare(
		mq.name, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err != nil {
		mq.logger.Error("error declaring queue", "error", err.Error())
		return err
	}

	return nil
}

func (mq *RabbitMQClient) Publish(data []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := mq.ch.PublishWithContext(ctx,
		"",            // exchange
		mq.queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		},
	)

	if err != nil {
		mq.logger.Error("error publishing message", "error", err.Error())
		return err
	}

	return nil
}

func (mq *RabbitMQClient) Consume(handler messagehandler.IMessageHandler) {

	msgs, err := mq.ch.Consume(
		mq.queue.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)

	if err != nil {
		mq.logger.Error("error consuming message", "error", err.Error())
		return
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			mq.logger.Info("received message", "message", string(d.Body))
			handler.HandleMessage(d.Body)
		}
	}()

	<-forever
}
