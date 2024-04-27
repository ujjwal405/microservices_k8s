package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Fetch interface {
	GetData(token string) (UserData, error)
}
type Middleware struct {
	RPCService Fetch
}

func NewMiddleware(svc Fetch) *Middleware {
	return &Middleware{
		RPCService: svc,
	}
}
func (md *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"err": "No token"})
			c.Abort()
			return
		}
		userdata, err := md.RPCService.GetData(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
			c.Abort()
			return
		}
		if !userdata.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"err": "Please login again"})
			c.Abort()
			return
		}
		c.Set("userid", userdata.Userid)
		c.Set("username", userdata.Username)
		c.Next()
	}
}
