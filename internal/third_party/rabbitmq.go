package third_party

import (
	"chat-apps/internal/repository"
	"chat-apps/internal/util"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Connection *amqp.Connection
	Channel    repository.RabbitMQChannel
	Queue      amqp.Queue
}

func NewRabbitMQ() (*RabbitMQ, error) {
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	dial := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)

	conn, err := amqp.Dial(dial)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"broadcast_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQ{
		Connection: conn,
		Channel:    ch,
		Queue:      q,
	}, nil
}

func (r *RabbitMQ) StartConsumer(queueName string, worker *util.NotificationWorker) error {
	msgs, err := r.Channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			worker.ProcessTask(d)
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() error {
	var err error
	if r.Channel != nil {
		err = r.Channel.Close()
		if err != nil {
			return err
		}
	}

	if r.Connection != nil {
		err = r.Connection.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
