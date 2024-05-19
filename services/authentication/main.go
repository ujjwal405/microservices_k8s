package main

import (
	"log"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	grpcserver "github.com/Ujjwal405/microservices/services/authentication/grpc"
	user "github.com/Ujjwal405/microservices/services/authentication/json"
	"github.com/Ujjwal405/microservices/services/authentication/json/router"

	"github.com/Ujjwal405/microservices/services/authentication/commonpkg/queue"
)

func main() {
	grpcadd := os.Getenv("ADDRESS_AUTH")
	rabbitmqadd := os.Getenv("RABBITMQ_ADD")
	go func() {
		err := grpcserver.RunGRPCServer(grpcadd)
		if err != nil {
			log.Panic(err)
		}
	}()
	mongoclient := user.Instancedb()
	mongoconn := user.OpenCollection(mongoclient, "microservice")
	Database := user.NewDatabase(mongoconn)
	amqpcon, amqpch, err := queue.Producer(rabbitmqadd)
	if err != nil {
		log.Panic(err)
	}
	Queue := user.NewQueue(amqpcon, amqpch)
	Memory := user.NewMemory()
	Service := user.NewService(Database, Memory, Queue)
	Handler := user.NewHandler(Service)
	router.InitRoutes(Handler)
	go func() {
		router.Start("0.0.0.0:8080")
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	amqpch.Close()
	amqpcon.Close()
	os.Exit(0)
}
