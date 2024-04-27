package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

var validate = validator.New()

type Handler struct {
	userservice *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		userservice: s,
	}
}
func (h *Handler) Signup(c *gin.Context) {
	var user_signup User
	if err := c.BindJSON(&user_signup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if err := validate.Struct(user_signup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	if err := h.userservice.Usersignup(user_signup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "verify the email"})

}
func (h *Handler) Login(c *gin.Context) {
	var user_login User
	if err := c.BindJSON(&user_login); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if err := validate.Struct(user_login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	token, err := h.userservice.UserLogin(user_login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.Header("token", token)
	c.JSON(http.StatusOK, gin.H{"success": "login successfull"})

}
func (h *Handler) Check(c *gin.Context) {
	var usrcode UserCode
	if err := c.BindJSON(&usrcode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	if usrcode.Usercode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Please provide the code"})
		return
	}
	if err := h.userservice.CheckUser(usrcode.Usercode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "successfully signed up"})
}
