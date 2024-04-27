package user

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func MainHandler(MQ MockQueue, MC MockCache, MR MockRepo) *Handler {
	svc := NewService(MR, MC, MQ)
	handler := NewHandler(svc)
	return handler
}
func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
