package repository

import (
	"github.com/streadway/amqp"
)

// Define an interface that includes only the methods your service needs
type RabbitMQChannel interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	Close() error
}
