package main

import (
	"log"
	"os"

	product "github.com/Ujjwal405/microservices/services/product"
	"github.com/Ujjwal405/microservices/services/product/commonpkg/middleware"
	"github.com/Ujjwal405/microservices/services/product/router"
)

func main() {
	mongoclient := product.Instancedb()
	productcollection := product.OpenCollection(mongoclient, "product")
	database := product.NewProductDatabase(productcollection)
	handler := product.NewHandler(database)
	authadd := os.Getenv("AUTHENTICATION_GRPC_ADD")
	authrpcclient, err := middleware.NewRpcClient(authadd)
	if err != nil {
		log.Fatal(err)
	}

	mid := middleware.NewMiddleware(middleware.NewClient(authrpcclient))
	router.InitRoutes(handler, mid)
	log.Fatal(router.Start("0.0.0.0:8080"))
}
