package grpc

import (
	"context"

	"github.com/pete-robinson/set-maker-grpc/internal/grpc/api"
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/domain"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Server) GetSong(ctx context.Context, id *wrapperspb.StringValue) (*domain.Song, error) {
	return &domain.Song{}, nil
}

func (s *Server) CreateSong(ctx context.Context, req *api.CreateSongRequest) (*domain.Song, error) {
	return &domain.Song{}, nil
}

func (s *Server) UpdateSong(ctx context.Context, req *api.UpdateSongRequest) (*domain.Song, error) {
	return &domain.Song{}, nil
}

func (s *Server) DeleteSong(ctx context.Context, id *wrapperspb.StringValue) (*api.DeleteSongResponse, error) {
	return &api.DeleteSongResponse{}, nil
}

func (s *Server) ListSongs(ctx context.Context, req *api.ListSongsRequest) (*api.ListSongsResponse, error) {
	return &api.ListSongsResponse{}, nil
}
