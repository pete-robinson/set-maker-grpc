package service

import (
	"github.com/google/uuid"
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/domain"
)

func (s *Service) CreateArtist(artist *domain.Artist) (*domain.Artist, error) {
	artist.Uuid = uuid.New().String()
	return artist, nil
}
