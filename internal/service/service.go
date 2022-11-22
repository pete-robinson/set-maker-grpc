package service

import (
	"context"

	"github.com/google/uuid"
	repository "github.com/pete-robinson/set-maker-grpc/internal/repository/ddb"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
)

type Repository interface {
	ListArtists(context.Context, int32, string) (*repository.ArtistList, error)
	GetArtist(context.Context, uuid.UUID) (*setmakerpb.Artist, error)
	PutArtist(context.Context, *setmakerpb.Artist) error
	DeleteArtist(context.Context, uuid.UUID) error

	ListSongs(context.Context, int32, string) (*repository.SongList, error)
	ListSongsByArtist(context.Context, int32, string, string) (*repository.SongList, error)
	GetSong(context.Context, uuid.UUID) (*setmakerpb.Song, error)
	PutSong(context.Context, *setmakerpb.Song) error
	DeleteSong(context.Context, uuid.UUID) error
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
