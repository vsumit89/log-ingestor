package message_queue

import (
	"logswift/internal/app/config"
	messagehandler "logswift/internal/message_queue/message_handler"
	"logswift/internal/message_queue/rabbitmq"
)

type IMessageQueue interface {
	Connect(cfg config.MQConfig) error
	DeclareQueue() error
	Publish([]byte) error
	Consume(messagehandler.IMessageHandler)
}

func NewMessageQueue() IMessageQueue {
	return rabbitmq.NewRabbitMQClient()
}
