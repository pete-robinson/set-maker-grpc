package service

import (
	"github.com/google/uuid"
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/domain"
)

type Repository interface {
	ListArtists() ([]*domain.Artist, error)
	GetArtist(uuid.UUID) (*domain.Artist, error)
	CreateArtist(*domain.Artist) (*domain.Artist, error)
	UpdateArtist(*domain.Artist) (*domain.Artist, error)
	DeleteArtist(uuid.UUID) error
}

type Service struct {
	repository Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}
