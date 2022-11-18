package service

import (
	"context"

	"github.com/google/uuid"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
)

type Repository interface {
	ListArtists(context.Context) ([]*setmakerpb.Artist, error)
	GetArtist(context.Context, uuid.UUID) (*setmakerpb.Artist, error)
	PutArtist(context.Context, *setmakerpb.Artist) error
	DeleteArtist(context.Context, uuid.UUID) error
}

type Notifier interface {
	RaiseArtistCreatedEvent(context.Context, *setmakerpb.Artist) error
}

type Service struct {
	repository Repository
	snsClient  Notifier
}

func NewService(repo Repository, sns Notifier) *Service {
	return &Service{
		repository: repo,
		snsClient:  sns,
	}
}
