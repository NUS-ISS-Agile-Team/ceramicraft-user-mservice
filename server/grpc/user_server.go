package grpc

import (
	"context"

	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/common/userpb"
	"github.com/NUS-ISS-Agile-Team/ceramicraft-user-mservice/server/log"
)

type UserServervice struct {
	userpb.UnimplementedUserServiceServer
}

func (s *UserServervice) SayHello(ctx context.Context, in *userpb.HelloRequest) (*userpb.HelloResponse, error) {
	log.Logger.Infof("Received: %v", in.GetName())
	return &userpb.HelloResponse{Message: "Hello " + in.GetName()}, nil
}
