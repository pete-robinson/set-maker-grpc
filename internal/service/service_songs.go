package service

import (
	"context"

	"github.com/google/uuid"
	repository "github.com/pete-robinson/set-maker-grpc/internal/repository/ddb"
	"github.com/pete-robinson/set-maker-grpc/internal/utils"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)


func (s *Service) ListSongs(ctx context.Context, limit int32, cursor string) (*repository.SongList, error) {
	res, err := s.repository.ListSongs(ctx, limit, cursor)
	if err != nil {
		return nil, err
	}

	logger.WithFields(logger.Fields{
		"result count": res.Count,
		"cursor": res.Cursor,
	}).Info("Results found")

	return res, nil
}


func (s *Service) ListSongsByArtist(ctx context.Context, limit int32, cursor string, artistId string) (*repository.SongList, error) {
	// artistID to uuid
	a, err := uuid.Parse(artistId);
	if err != nil {
		logger.WithField("uuid", artistId).Errorf("Could not parse artist UUID: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid artist Id")
	}

	if _, err := s.repository.GetArtist(ctx, a); err != nil {
		logger.WithField("artistId", artistId).Errorf("Error fetching artist: %s", err)
		return nil, status.Error(codes.NotFound, "Unknown artist")
	}

	res, err := s.repository.ListSongsByArtist(ctx, limit, cursor, artistId)
	if err != nil {
		return nil, err
	}

	logger.WithFields(logger.Fields{
		"result count": res.Count,
		"cursor": res.Cursor,
		"artist": artistId,
	}).Info("Results found")

	return res, nil
}


func (s *Service) GetSong(ctx context.Context, id uuid.UUID) (*setmakerpb.Song, error) {
	song, err := s.repository.GetSong(ctx, id)
	if err != nil {
		logger.WithField("id", id).Errorf("Could not fetch song: %s", err)
		return nil, err
	}

	return song, nil
}


func (s *Service) CreateSong(ctx context.Context, song *setmakerpb.Song) (*setmakerpb.Song, error) {
	// artistID to uuid
	artistId, err := uuid.Parse(song.ArtistId);
	if err != nil {
		logger.WithField("uuid", song.ArtistId).Errorf("Could not parse artist UUID: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid artist Id")
	}

	// validate the artist exists
	if _, err := s.GetArtist(ctx, artistId); err != nil {
		logger.WithField("artistId", artistId).Errorf("Error locating artist for song: %s", err)
		return nil, err
	}

	// init UUID and meta
	song.Id = uuid.New().String()
	song.Metadata = &setmakerpb.Metadata{}
	utils.SetMetaData(song.Metadata)

	if err := s.repository.PutSong(ctx, song); err != nil {
		logger.WithField("data", song).Errorf("Could not create song: %s", err)
		return nil, err
	}

	return song, nil
}


func (s *Service) UpdateSong(ctx context.Context, song *setmakerpb.Song) (*setmakerpb.Song, error) {
	// artistID to uuid
	artistId, err := uuid.Parse(song.ArtistId);
	if err != nil {
		logger.WithField("uuid", song.ArtistId).Errorf("Could not parse artist UUID: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid artist Id")
	}

	// validate the artist exists
	if _, err := s.GetArtist(ctx, artistId); err != nil {
		logger.WithField("artistId", artistId).Errorf("Error locating artist for song: %s", err)
		return nil, err
	}

	// fetch the song to update
	targetId, err := uuid.Parse(song.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Error fetching song to update")
	}

	target, err := s.GetSong(ctx, targetId)
	if target == nil {
		return nil, status.Error(codes.NotFound, "Song to update does not exist")
	}

	song.Metadata = target.Metadata
	utils.SetMetaData(song.Metadata)

	// update artist
	if err = s.repository.PutSong(ctx, song); err != nil {
		return nil, err
	}

	return song, nil
}


func (s *Service) DeleteSong(ctx context.Context, id uuid.UUID) error {
	if err := s.repository.DeleteSong(ctx, id); err != nil {
		return err
	}

	return nil
}
