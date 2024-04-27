package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `bson:"_id"`
	Email    string             `json:"email" bson:"email" validate:"email,required"`
	Username string             `json:"username" bson:"username" validate:"required,min=4,max=8"`
	Password string             `json:"password" bson:"password" validate:"required,min=5,max=8"`
	User_id  string             `bson:"user_id"`
}
type Details struct {
	TTL    time.Time
	Detail User
}
type Mail struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type UserCode struct {
	Usercode string `json:"code"`
}
