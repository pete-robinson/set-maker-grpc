package grpc

import (
	"context"

	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Server) GetSong(ctx context.Context, id *wrapperspb.StringValue) (*setmakerpb.Song, error) {
	return &setmakerpb.Song{}, nil
}

func (s *Server) CreateSong(ctx context.Context, req *setmakerpb.CreateSongRequest) (*setmakerpb.Song, error) {
	return &setmakerpb.Song{}, nil
}

func (s *Server) UpdateSong(ctx context.Context, req *setmakerpb.UpdateSongRequest) (*setmakerpb.Song, error) {
	return &setmakerpb.Song{}, nil
}

func (s *Server) DeleteSong(ctx context.Context, id *wrapperspb.StringValue) (*setmakerpb.DeleteSongResponse, error) {
	return &setmakerpb.DeleteSongResponse{}, nil
}

func (s *Server) ListSongs(ctx context.Context, req *setmakerpb.ListSongsRequest) (*setmakerpb.ListSongsResponse, error) {
	return &setmakerpb.ListSongsResponse{}, nil
}
