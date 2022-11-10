package grpc

import (
	"context"
	"time"

	"github.com/pete-robinson/set-maker-grpc/internal/grpc/api"
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/domain"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Server) GetArtist(ctx context.Context, id *wrapperspb.StringValue) (*domain.Artist, error) {
	return &domain.Artist{
		Uuid:  "fdsfsfs",
		Name:  "test",
		Image: "image",
		Metadata: &domain.Metadata{
			CreatedAt: time.Now().UTC().String(),
			UpdatedAt: time.Now().UTC().String(),
		},
	}, nil
}

func (s *Server) CreateArtist(ctx context.Context, req *api.CreateArtistRequest) (*domain.Artist, error) {
	return &domain.Artist{}, nil
}

func (s *Server) UpdateArtist(ctx context.Context, artist *domain.Artist) (*domain.Artist, error) {
	return &domain.Artist{}, nil
}

func (s *Server) DeleteArtist(ctx context.Context, id *wrapperspb.StringValue) (*api.DeleteArtistResponse, error) {
	return &api.DeleteArtistResponse{}, nil
}

func (s *Server) ListArtists(ctx context.Context, req *api.ListArtistsRequest) (*api.ListArtistsResponse, error) {
	return &api.ListArtistsResponse{}, nil
}
