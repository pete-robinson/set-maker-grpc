package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/domain"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const TableName = "artists"

func (d *DynamoRepository) ListArtists(ctx context.Context) ([]*domain.Artist, error) {
	// create artist
	return []*domain.Artist{}, nil
}

func (d *DynamoRepository) GetArtist(ctx context.Context, id uuid.UUID) (*domain.Artist, error) {
	// create key map
	keys, err := attributevalue.MarshalMap(map[string]string{
		"Id": *aws.String(id.String()),
	})
	if err != nil {
		logger.WithField("id", id).Errorf("Could not marshalmap: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid UUID")
	}

	// fetch item from dynamo
	data, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       keys,
	})
	if err != nil {
		logger.WithField("id", id).Errorf("Error fetching from dynamo: %s", err)
		return nil, status.Error(codes.Internal, "Error fetching result")
	}

	// check an item was returned
	if data.Item == nil {
		logger.WithField("id", id).Error("No artist found for ID")
		return nil, status.Error(codes.NotFound, "Artist not found")
	}

	// fetch was successful
	logger.WithField("data", data.Item).Info("Artist found")

	// unmarshal response
	res := &domain.Artist{}
	err = attributevalue.UnmarshalMap(data.Item, res)
	if err != nil {
		logger.WithField("data", data.Item).Errorf("could not unmarshal item: %s", err)
		return nil, status.Error(codes.Internal, "Error unmarshaling artist data")
	}

	return res, nil
}

func (d *DynamoRepository) PutArtist(ctx context.Context, artist *domain.Artist) error {
	// create attribute value map
	item, err := attributevalue.MarshalMap(artist)
	if err != nil {
		logger.WithField("data", artist).Errorf("Could not marshalmap: %s", err)
		return status.Error(codes.InvalidArgument, "Could not map input values for artist")
	}

	// PutItem to dynamo
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})
	if err != nil {
		logger.WithField("data", artist).Errorf("Could not PutItem: %s", err)
		return status.Error(codes.Internal, "Failed to persist artist")
	}

	logger.WithField("id", artist.Id).Info("Artist persisted successfully")
	return nil
}

func (d *DynamoRepository) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	logger.WithField("id", id).Infof("Deleting artist")

	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		logger.WithField("id", id).Errorf("Could not delete artist: %s", err)
		return status.Error(codes.Internal, "Artist could not be deleted")
	}

	return nil
}
