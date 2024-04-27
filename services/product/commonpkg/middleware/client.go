package middleware

import (
	"context"
	"errors"
	"time"

	"github.com/Ujjwal405/microservices/services/product/commonpkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func NewRpcClient(address string) (proto.AuthServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := proto.NewAuthServiceClient(conn)
	return c, nil
}

type UserData struct {
	Userid   string
	Valid    bool
	Username string
}
type Client struct {
	client proto.AuthServiceClient
}

func NewClient(c proto.AuthServiceClient) *Client {
	return &Client{
		client: c,
	}
}
func (c *Client) GetData(token string) (UserData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &proto.Request{
		Token: token,
	}
	res, err := c.client.FetchData(ctx, req)
	if err != nil {
		if s, ok := status.FromError(err); ok {
			if s.Code() == codes.DeadlineExceeded {
				return UserData{}, errors.New("deadline exceeded")
			} else if s.Code() == codes.Canceled {
				return UserData{}, errors.New("context cancelled")
			} else {
				return UserData{}, errors.New(" Internal server error")
			}
		}
	}
	usrdata := UserData{
		Userid:   res.GetUserid(),
		Valid:    res.GetValid(),
		Username: res.GetUsername(),
	}
	return usrdata, nil
}
