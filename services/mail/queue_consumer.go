package mail

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Mailing interface {
	SendGmail(code string, to string) error
}

type Consumer struct {
	conn      *amqp.Connection
	deliver   <-chan amqp.Delivery
	ch        *amqp.Channel
	mailer    Mailing
	semaphore chan struct{}
}

func NewConsumer(conn *amqp.Connection, delv <-chan amqp.Delivery, chann *amqp.Channel, ml Mailing) *Consumer {
	return &Consumer{
		conn:      conn,
		deliver:   delv,
		ch:        chann,
		mailer:    ml,
		semaphore: make(chan struct{}, 3),
	}
}
func (cn *Consumer) ConsumeFromQueue() {
	cn.Worker()
}
func (cn *Consumer) Worker() {
	for msg := range cn.deliver {
		cn.semaphore <- struct{}{}
		go func(msg amqp.Delivery) {
			defer func() {
				<-cn.semaphore
			}()
			var body Mailstruct
			if err := json.Unmarshal(msg.Body, &body); err != nil {
				log.Println(err.Error())
			}

			if err := cn.mailer.SendGmail(body.Code, body.Email); err != nil {

				log.Println(err)
				msg.Nack(false, false)
				//continue
			}
			msg.Ack(true)
		}(msg)
	}

}
func (cn *Consumer) Close() {
	cn.ch.Close()
	cn.conn.Close()
}
