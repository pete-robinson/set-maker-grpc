package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/api"
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/domain"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *Server) GetArtist(ctx context.Context, id *wrapperspb.StringValue) (*domain.Artist, error) {
	logger.WithField("id", id).Info("Attempting to fetch artist")

	// parse UUID
	uuid, err := uuid.Parse(id.Value)
	if err != nil {
		logger.WithField("uuid", id.Value).Errorf("Could not parse UUID: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid artist Id")
	}

	// fetch artist from service
	artist, err := s.service.GetArtist(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (s *Server) CreateArtist(ctx context.Context, req *api.CreateArtistRequest) (*domain.Artist, error) {
	logger.WithField("request", req).Info("Creating artist")

	artist := &domain.Artist{
		Name:  req.Name,
		Image: req.Image,
	}

	created, err := s.service.CreateArtist(ctx, artist)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Server) UpdateArtist(ctx context.Context, req *api.UpdateArtistRequest) (*domain.Artist, error) {
	// validate the UUID
	_, err := uuid.Parse(req.Id)
	if err != nil {
		logger.WithField("uuid", req.Id).Errorf("Could not parse UUID: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid artist Id")
	}

	artist := &domain.Artist{
		Id:    req.Id,
		Name:  req.Name,
		Image: req.Image,
	}

	// attempt update
	resp, err := s.service.UpdateArtist(ctx, artist)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Server) DeleteArtist(ctx context.Context, id *wrapperspb.StringValue) (*api.DeleteArtistResponse, error) {
	logger.WithField("id", id.GetValue()).Info("Deleting artist")

	// parse UUID
	uuid, err := uuid.Parse(id.GetValue())
	if err != nil {
		logger.WithField("id", id.GetValue()).Errorf("ID will not parse %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid data for Id")
	}

	logger.WithField("uuid", uuid).Info("Generated UUID")

	// run delete
	resp := &api.DeleteArtistResponse{
		Id:      uuid.String(),
		Deleted: false,
	}

	err = s.service.DeleteArtist(ctx, uuid)
	if err != nil {
		return nil, err
	}

	logger.WithField("id", uuid.String()).Info("Artist deleted successfully")
	resp.Deleted = true

	return resp, nil
}

func (s *Server) ListArtists(ctx context.Context, req *api.ListArtistsRequest) (*api.ListArtistsResponse, error) {
	return &api.ListArtistsResponse{}, nil
}
