package grpcserver

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Signedetails struct {
	User_id   string
	User_name string
	jwt.StandardClaims
}
type Authservice struct {
}

func NewAuthService() *Authservice {
	return &Authservice{}
}
func (auth *Authservice) ValidateToken(SignedToken string) (Signedetails, error) {

	var SECRET_KEY string = os.Getenv("SECRET_KEY")
	token, err := jwt.ParseWithClaims(
		SignedToken,
		&Signedetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return Signedetails{}, err
	}

	claims, ok := token.Claims.(*Signedetails)
	if !ok {

		err = errors.New("token invalid")
		return Signedetails{}, err
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return Signedetails{}, err
	}
	return *claims, err
}
