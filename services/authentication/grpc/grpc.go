package grpcserver

import (
	"context"
	"net"

	"github.com/Ujjwal405/microservices/services/authentication/commonpkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	expire  string = "token expired"
	invalid string = "token invalid"
)

type Validate interface {
	ValidateToken(SignedToken string) (Signedetails, error)
}
type GRPCAuthServer struct {
	authservice Validate
	proto.UnimplementedAuthServiceServer
}

func RunGRPCServer(address string) error {
	authsvc := NewAuthService()
	newGrpcAuthService := NewGRPCAuthServer(authsvc)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{}
	server := grpc.NewServer(opts...)
	proto.RegisterAuthServiceServer(server, newGrpcAuthService)
	return server.Serve(ln)

}
func NewGRPCAuthServer(authsvc Validate) *GRPCAuthServer {
	return &GRPCAuthServer{
		authservice: authsvc,
	}
}
func (server *GRPCAuthServer) FetchData(ctx context.Context, req *proto.Request) (*proto.Response, error) {
	var Isok bool = true
	uid := req.GetToken()
	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "context cancelled")
	}
	claim, err := server.authservice.ValidateToken(uid)
	if err != nil {
		if err.Error() == expire || err.Error() == invalid {
			Isok = false
		} else {
			return nil, status.Errorf(codes.Internal, "%v", err)
		}
	}
	if ctx.Err() == context.DeadlineExceeded {
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	}
	res := &proto.Response{
		Userid:   claim.User_id,
		Valid:    Isok,
		Username: claim.User_name,
	}
	return res, nil
}
