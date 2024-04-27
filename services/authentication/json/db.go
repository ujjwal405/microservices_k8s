package user

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database struct {
	Usercollection *mongo.Collection
}

func NewDatabase(usecollection *mongo.Collection) *Database {
	return &Database{
		Usercollection: usecollection,
	}
}
func (db *Database) CountDocuments(email string) (bool, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	count, err := db.Usercollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, errors.New("this email already exists")
	}
	return false, nil
}
func (db *Database) InsertUser(user User) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	_, inserterr := db.Usercollection.InsertOne(ctx, user)
	if inserterr != nil {
		return inserterr
	}
	return nil
}
func (db *Database) FetchUser(email string) (User, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	isok, err := db.CountDocuments(email)
	if !isok {
		if err != nil {
			return User{}, errors.New("internal server error")
		} else {
			return User{}, errors.New("this email doesn't exist")
		}

	}
	var founduser User
	err = db.Usercollection.FindOne(ctx, bson.M{"email": email}).Decode(&founduser)
	if err != nil {
		return User{}, err
	}

	return founduser, nil
}
