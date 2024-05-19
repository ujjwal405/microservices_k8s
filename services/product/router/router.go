package router

import (
	product "github.com/Ujjwal405/microservices/services/product"
	"github.com/Ujjwal405/microservices/services/product/commonpkg/middleware"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRoutes(handler *product.Handler, mid *middleware.Middleware) {
	r = gin.Default()
	r.GET("/health", handler.Health)
	r.Use(mid.AuthMiddleware())
	r.POST("/addproduct", handler.AddProduct)
	r.GET("/searchproduct", handler.SearchProduct)

}
func Start(address string) error {
	return r.Run(address)
}
