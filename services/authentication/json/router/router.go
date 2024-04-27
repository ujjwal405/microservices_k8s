package router

import (
	user "github.com/Ujjwal405/microservices/services/authentication/json"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func InitRoutes(handler *user.Handler) {
	r = gin.Default()
	r.POST("/signup", handler.Signup)
	r.POST("/login", handler.Login)
	r.POST("/check", handler.Check)
}
func Start(address string) error {
	return r.Run(address)
}
