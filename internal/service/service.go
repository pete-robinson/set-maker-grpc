package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/domain"
)

type Repository interface {
	ListArtists(context.Context) ([]*domain.Artist, error)
	GetArtist(context.Context, uuid.UUID) (*domain.Artist, error)
	PutArtist(context.Context, *domain.Artist) error
	DeleteArtist(context.Context, uuid.UUID) error
}

type Service struct {
	repository Repository
	snsClient  *sns.Client
}

func NewService(repo Repository, sns *sns.Client) *Service {
	return &Service{
		repository: repo,
		snsClient:  sns,
	}
}
