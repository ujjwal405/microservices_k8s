package queue

import (
	"github.com/streadway/amqp"
)

func Producer(url string) (*amqp.Connection, *amqp.Channel, error) {
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, nil, err
	}

	// declaring queue with its properties over the the channel opened
	_, err = channel.QueueDeclare(
		"testing", // name
		false,     // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		return nil, nil, err

	}

	return connection, channel, nil
}
