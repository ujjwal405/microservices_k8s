package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	mail "github.com/Ujjwal405/microservices/services/mail"
	"github.com/Ujjwal405/microservices/services/mail/commonpkg/queue"
)

func main() {
	rabbitmqadd := os.Getenv("RABBITMQ_ADD")
	frompassEmail := os.Getenv("FROM_EMAIL_ADDRESS")
	fromEmailPassword := os.Getenv("FROM_PASSWORD_EMAIL")
	name := os.Getenv("NAME")
	mailer := mail.NewMailer(name, frompassEmail, fromEmailPassword)
	qconn, chdelivery, qch, err := queue.Consumer(rabbitmqadd)
	if err != nil {
		panic(err)
	}
	consumer := mail.NewConsumer(qconn, chdelivery, qch, mailer)
	go func() {
		consumer.ConsumeFromQueue()
	}()
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic(err)
		}
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	consumer.Close()
	os.Exit(0)
}
