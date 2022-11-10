package grpc

import (
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/api"
	"github.com/pete-robinson/set-maker-grpc/internal/service"
)

type Server struct {
	service *service.Service
	api.UnimplementedSetMakerServiceServer
}

func NewServer(svc *service.Service) (*Server, error) {
	return &Server{
		service: svc,
	}, nil
}
