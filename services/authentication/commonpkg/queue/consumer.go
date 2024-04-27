package queue

import "github.com/streadway/amqp"

func Consumer(url string) (*amqp.Connection, <-chan amqp.Delivery, *amqp.Channel, error) {
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, nil, err
	}
	channel, err := connection.Channel()
	if err != nil {
		return nil, nil, nil, err
	}

	// declaring queue with its properties over the the channel opened
	msgs, err := channel.Consume(
		"testing", // queue
		"",        // consumer
		false,     // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       //args
	)
	if err != nil {
		return nil, nil, nil, err

	}
	return connection, msgs, channel, nil

}
