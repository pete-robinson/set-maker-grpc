package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
)

type Repository interface {
	ListArtists(context.Context) ([]*setmakerpb.Artist, error)
	GetArtist(context.Context, uuid.UUID) (*setmakerpb.Artist, error)
	PutArtist(context.Context, *setmakerpb.Artist) error
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
