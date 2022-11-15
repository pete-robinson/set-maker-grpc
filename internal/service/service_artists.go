package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/google/uuid"
	"github.com/pete-robinson/set-maker-grpc/internal/utils"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const SnsTopic = "create-artist.fifo"

func (s *Service) GetArtist(ctx context.Context, id uuid.UUID) (*setmakerpb.Artist, error) {
	artist, err := s.repository.GetArtist(ctx, id)
	if err != nil {
		logger.WithField("id", id).Errorf("Could not fetch artist: %s", err)
		return nil, err
	}

	return artist, nil
}

func (s *Service) CreateArtist(ctx context.Context, artist *setmakerpb.Artist) (*setmakerpb.Artist, error) {
	// init UUID and meta
	artist.Id = uuid.New().String()
	utils.SetMetaData(&setmakerpb.Metadata{})

	err := s.repository.PutArtist(ctx, artist)
	if err != nil {
		logger.WithField("data", artist).Errorf("Could not create artist: %s", err)
		return nil, err
	}

	notificationTopic := "arn:aws:sns:eu-west-1:194252907131:create-artist"

	// publish sns message
	snsIn := &sns.PublishInput{
		Message:  &artist.Id,
		TopicArn: &notificationTopic,
	}

	res, err := s.snsClient.Publish(ctx, snsIn)
	if err != nil {
		logger.WithField("id", artist.Id).Errorf("Could not publish to SNS topic: %s", err)
	}

	logger.WithFields(logger.Fields{
		"id":        artist.Id,
		"messageId": *res.MessageId,
	}).Info("Published new artist to SNS topic")

	return artist, nil
}

func (s *Service) UpdateArtist(ctx context.Context, artist *setmakerpb.Artist) (*setmakerpb.Artist, error) {
	// fetch the artist to update
	targetId, err := uuid.Parse(artist.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Error fetching artist to update")
	}

	target, err := s.GetArtist(ctx, targetId)
	if target == nil {
		return nil, status.Error(codes.NotFound, "Artist to update does not exist")
	}

	// reset the data
	target.Name = artist.Name
	target.Image = artist.Image
	utils.SetMetaData(target.Metadata)

	// update artist
	err = s.repository.PutArtist(ctx, target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

func (s *Service) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	err := s.repository.DeleteArtist(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
