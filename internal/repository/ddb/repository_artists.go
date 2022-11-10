package repository

import (
	"github.com/google/uuid"
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/domain"
)

func (d *DynamoRepository) ListArtists() ([]*domain.Artist, error) {
	// create artist
	return []*domain.Artist{}, nil
}

func (d *DynamoRepository) GetArtist(id uuid.UUID) (*domain.Artist, error) {
	// create artist
	return &domain.Artist{}, nil
}

func (d *DynamoRepository) CreateArtist(artist *domain.Artist) (*domain.Artist, error) {
	// create artist
	return artist, nil
}

func (d *DynamoRepository) UpdateArtist(artist *domain.Artist) (*domain.Artist, error) {
	// create artist
	return artist, nil
}

func (d *DynamoRepository) DeleteArtist(id uuid.UUID) error {
	// create artist
	return nil
}
