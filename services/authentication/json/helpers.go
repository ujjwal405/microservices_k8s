package user

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Signedetails struct {
	User_id   string
	User_name string
	jwt.StandardClaims
}

func GenerateRandom() string {
	uuid := uuid.New().String()
	truncated := uuid[:4]
	return truncated
}
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	check := true
	msg := ""
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))
	if err != nil {
		msg = "password incorect"
		check = false
	}
	return check, msg
}
func HashPassword(userpassword string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(userpassword), 14)
	if err != nil {
		return ""
	}
	return string(bytes)
}
func GenerateAlltoken(user User) (string, error) {
	//err := config.Loadconfig("../../..")
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return "", err
	//}

	var SECRET_KEY string = os.Getenv("SECRET_KEY")
	claims := &Signedetails{
		User_id:   user.User_id,
		User_name: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {

		return "", err
	}
	log.Println(token)
	return token, nil
}
