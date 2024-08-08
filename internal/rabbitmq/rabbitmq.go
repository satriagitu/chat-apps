package rabbitmq

import (
	"chat-apps/internal/util"
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

func InitRabbitMQ() (*amqp.Connection, *amqp.Channel, amqp.Queue, error) {
	user := os.Getenv("RABBITMQ_USER")
	pass := os.Getenv("RABBITMQ_PASS")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	dial := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	// dial := "amqp://guest:guest@rabbitmq:5672/"
	conn, err := amqp.Dial(dial)
	if err != nil {
		return nil, nil, amqp.Queue{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, nil, amqp.Queue{}, err
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
		return nil, nil, amqp.Queue{}, err
	}

	return conn, ch, q, nil
}

func StartConsumer(ch *amqp.Channel, queueName string, worker *util.NotificationWorker) error {
	msgs, err := ch.Consume(
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
