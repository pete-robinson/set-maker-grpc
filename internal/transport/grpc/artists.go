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

func (s *Server) GetArtist(ctx context.Context, id *wrapperspb.StringValue) (*setmakerpb.Artist, error) {
	logger.WithField("id", id).Info("GRPC: Fetching artist")

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

func (s *Server) CreateArtist(ctx context.Context, req *setmakerpb.CreateArtistRequest) (*setmakerpb.Artist, error) {
	logger.WithField("request", req).Info("GRPC: Creating artist")

	artist := &setmakerpb.Artist{
		Name:  req.Name,
		Image: req.Image,
	}

	created, err := s.service.CreateArtist(ctx, artist)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Server) UpdateArtist(ctx context.Context, req *setmakerpb.UpdateArtistRequest) (*setmakerpb.Artist, error) {
	logger.WithField("request", req).Info("GRPC: Updating Artist")

	// validate the UUID
	if _, err := uuid.Parse(req.Id); err != nil {
		logger.WithField("uuid", req.Id).Errorf("Could not parse UUID: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid artist Id")
	}

	artist := &setmakerpb.Artist{
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

func (s *Server) DeleteArtist(ctx context.Context, id *wrapperspb.StringValue) (*setmakerpb.DeleteArtistResponse, error) {
	logger.WithField("id", id.GetValue()).Info("GRPC: Deleting artist")

	// parse UUID
	uuid, err := uuid.Parse(id.GetValue())
	if err != nil {
		logger.WithField("id", id.GetValue()).Errorf("ID will not parse %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid data for Id")
	}

	logger.WithField("uuid", uuid).Info("Generated UUID")

	// run delete
	resp := &setmakerpb.DeleteArtistResponse{
		Id:      uuid.String(),
		Deleted: false,
	}

	if err = s.service.DeleteArtist(ctx, uuid); err != nil {
		return nil, err
	}

	logger.WithField("id", uuid.String()).Info("Artist deleted successfully")
	resp.Deleted = true

	return resp, nil
}

func (s *Server) ListArtists(ctx context.Context, req *setmakerpb.ListArtistsRequest) (*setmakerpb.ListArtistsResponse, error) {
	logger.WithField("req", req).Info("GRPC: Listing artists")

	resp, err := s.service.ListArtists(ctx, req.Limit, req.Cursor)
	if err != nil {
		logger.WithFields(logger.Fields{
			"limit": req.Limit,
			"cursor": req.Cursor,
		}).Errorf("Error listing artists: %s", err)
		return nil, err
	}

	r := &setmakerpb.ListArtistsResponse{
		Results: resp.Items,
		SearchAfter: resp.Cursor,
	}

	return r, nil
}
