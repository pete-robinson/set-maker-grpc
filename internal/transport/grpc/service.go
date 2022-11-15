package grpc

import (
	"github.com/pete-robinson/set-maker-grpc/internal/service"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
)

type Server struct {
	service *service.Service
	setmakerpb.UnimplementedSetMakerServiceServer
}

func NewServer(svc *service.Service) (*Server, error) {
	return &Server{
		service: svc,
	}, nil
}
