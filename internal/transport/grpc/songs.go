package grpc

import (
	"context"

	"github.com/google/uuid"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)


func (s *Server) GetSong(ctx context.Context, id *wrapperspb.StringValue) (*setmakerpb.Song, error) {
	logger.WithField("id", id).Info("GRPC: Fetching song")

	// parse UUID
	uuid, err := uuid.Parse(id.Value)
	if err != nil {
		logger.WithField("uuid", id.Value).Errorf("Could not parse UUID: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid song Id")
	}

	// fetch song from service
	song, err := s.service.GetSong(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return song, nil
}


func (s *Server) CreateSong(ctx context.Context, req *setmakerpb.CreateSongRequest) (*setmakerpb.Song, error) {
	logger.WithField("request", req).Info("GRPC: Creating song")

	song := &setmakerpb.Song{
		Title:  req.Title,
		ArtistId: req.ArtistId,
		Key: req.Key,
		Tonality: req.Tonality,
	}

	created, err := s.service.CreateSong(ctx, song)
	if err != nil {
		return nil, err
	}

	return created, nil
}


func (s *Server) UpdateSong(ctx context.Context, req *setmakerpb.UpdateSongRequest) (*setmakerpb.Song, error) {
	logger.WithField("request", req).Info("GRPC: Updating Song")

	// validate the UUID
	if _, err := uuid.Parse(req.Id); err != nil {
		logger.WithField("uuid", req.Id).Errorf("Could not parse UUID: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid song Id")
	}

	song := &setmakerpb.Song{
		Id: req.Id,
		Title:  req.Title,
		ArtistId: req.ArtistId,
		Key: req.Key,
		Tonality: req.Tonality,
	}

	// attempt update
	resp, err := s.service.UpdateSong(ctx, song)
	if err != nil {
		return nil, err
	}

	return resp, nil
}


func (s *Server) DeleteSong(ctx context.Context, id *wrapperspb.StringValue) (*setmakerpb.DeleteSongResponse, error) {
	logger.WithField("id", id.GetValue()).Info("GRPC: Deleting song")

	// parse UUID
	uuid, err := uuid.Parse(id.GetValue())
	if err != nil {
		logger.WithField("id", id.GetValue()).Errorf("ID will not parse %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid data for Id")
	}

	resp := &setmakerpb.DeleteSongResponse{
		Id:      uuid.String(),
		Deleted: false,
	}

	logger.WithField("uuid", uuid).Info("Generated UUID")

	// run delete
	if err = s.service.DeleteSong(ctx, uuid); err != nil {
		return resp, err
	}

	logger.WithField("id", uuid.String()).Info("Song deleted successfully")
	resp.Deleted = true

	return resp, nil
}


func (s *Server) ListSongs(ctx context.Context, req *setmakerpb.ListSongsRequest) (*setmakerpb.ListSongsResponse, error) {
	logger.WithField("req", req).Info("GRPC: Listing songs")

	resp, err := s.service.ListSongs(ctx, req.Limit, req.Cursor)
	if err != nil {
		logger.WithFields(logger.Fields{
			"limit": req.Limit,
			"cursor": req.Cursor,
		}).Errorf("Error listing songs: %s", err)
		return nil, err
	}

	r := &setmakerpb.ListSongsResponse{
		Results: resp.Items,
		SearchAfter: resp.Cursor,
	}

	return r, nil
}


func (s *Server) ListSongsByArtist(ctx context.Context, req *setmakerpb.ListSongsByArtistRequest) (*setmakerpb.ListSongsResponse, error) {
	logger.WithField("req", req).Info("GRPC: Listing songs")

	resp, err := s.service.ListSongsByArtist(ctx, req.Limit, req.Cursor, req.ArtistId)
	if err != nil {
		logger.WithFields(logger.Fields{
			"limit": req.Limit,
			"cursor": req.Cursor,
			"artistId": req.ArtistId,
		}).Errorf("Error listing songs by artist: %s", err)
		return nil, err
	}

	r := &setmakerpb.ListSongsResponse{
		Results: resp.Items,
		SearchAfter: resp.Cursor,
	}

	return r, nil
}
