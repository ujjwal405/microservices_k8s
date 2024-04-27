package user

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type Queue struct {
	Connection *amqp.Connection
	Wch        *amqp.Channel
}

func NewQueue(conn *amqp.Connection, ch *amqp.Channel) *Queue {
	return &Queue{
		Connection: conn,
		Wch:        ch,
	}

}
func (q *Queue) AddToQueue(mail Mail) error {
	payload, err := json.Marshal(mail)
	if err != nil {
		return err
	}
	err = q.Wch.Publish(
		"",
		"testing",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
	if err != nil {
		return err
	}
	return nil
}
