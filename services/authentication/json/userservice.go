package user

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	Repository Repo
	Inmemory   Cache
	Queue      queue
}

func NewService(repo Repo, cache Cache, q queue) *Service {
	return &Service{
		Repository: repo,
		Inmemory:   cache,
		Queue:      q,
	}
}
func (svc *Service) Usersignup(user User) error {
	isok, err := svc.Repository.CountDocuments(user.Email)
	if isok {
		return err
	}
	usertimeout := time.Now().Local().Add(35 * time.Second)
	detail := &Details{
		TTL:    usertimeout,
		Detail: user,
	}
	code := GenerateRandom()
	svc.Inmemory.AddClient(code, detail)
	svc.Inmemory.Awake()
	mail := Mail{
		Email: user.Email,
		Code:  code,
	}
	err = svc.Queue.AddToQueue(mail)
	if err != nil {
		return err
	}
	return nil

}
func (svc *Service) UserLogin(user User) (string, error) {
	isok, err := svc.Repository.CountDocuments(user.Email)
	if !isok {
		return "", err
	}
	dbuser, err := svc.Repository.FetchUser(user.Email)
	if err != nil {
		return "", err
	}
	isvalid, msg := VerifyPassword(dbuser.Password, user.Password)
	if !isvalid {
		return "", errors.New(msg)
	}
	token, err := GenerateAlltoken(dbuser)

	if err != nil {
		return "", err
	}
	return token, nil

}
func (svc *Service) CheckUser(code string) error {
	details, isok := svc.Inmemory.Check(code)
	if !isok {
		return errors.New("code doesn't exist")
	}
	hashedpass := HashPassword(details.Detail.Password)
	details.Detail.Password = hashedpass
	details.Detail.Id = primitive.NewObjectID()
	details.Detail.User_id = details.Detail.Id.Hex()
	if err := svc.Repository.InsertUser(details.Detail); err != nil {
		return err
	}
	return nil
}
