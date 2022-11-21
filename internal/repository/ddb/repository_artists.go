package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/pete-robinson/set-maker-grpc/internal/utils"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const TableName = "artists"

type ArtistList struct {
	Count int32
	Cursor string
	Items []*setmakerpb.Artist
}

func (d *DynamoRepository) ListArtists(ctx context.Context, limit int32, cursor string) (*ArtistList, error) {
	// decode the cursor
	c, err := utils.DecodeAttributeMap(cursor)
	if err != nil {
		logger.WithField("cursor", cursor).Errorf("ListArtists Repo: Could not decode cursor: %s", err)
		return nil, err
	}

	// build DDB scan input
	input := dynamodb.ScanInput{
		TableName: aws.String(TableName),
		Limit: &limit,
		ExclusiveStartKey: c,
	}

	logger.WithFields(logger.Fields{
		"limit": limit,
		"cursor": cursor,
	}).Info("ListArtists Repo: Scanning dynamo")

	// query DDB
	res, err := d.client.Scan(ctx, &input)
	if err != nil {
		logger.WithFields(logger.Fields{
			"limit": limit,
			"cursor": cursor,
			"input": input,
		}).Errorf("ListArtists Repo: Error scanning: %s", err)
	}

	logger.WithField("result count", res.Count).Info("ListArtists Repo: Results returned")

	// encode the return cursor
	returnCursor, err := utils.EncodeAttributeMap(res.LastEvaluatedKey)
	if err != nil {
		logger.WithField("lastEvaluatedKey", res.LastEvaluatedKey).Errorf("ListArtists Repo: Unable to encode attribute map: %s", err)
		return nil, err
	}

	// parse results
	var items []*setmakerpb.Artist
	if err = attributevalue.UnmarshalListOfMaps(res.Items, &items); err != nil {
		logger.Errorf("ListArtists Repo: Could not unmarshal results: %s", err)
		return nil, err
	}

	resp := &ArtistList{
		Count: res.Count,
		Cursor: string(returnCursor),
		Items: items,
	}

	return resp, nil
}

func (d *DynamoRepository) GetArtist(ctx context.Context, id uuid.UUID) (*setmakerpb.Artist, error) {
	// create key map
	keys, err := attributevalue.MarshalMap(map[string]string{
		"Id": *aws.String(id.String()),
	})
	if err != nil {
		logger.WithField("id", id).Errorf("GetArtists: Could not marshalmap: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid UUID")
	}

	// fetch item from dynamo
	data, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key:       keys,
	})
	if err != nil {
		logger.WithField("id", id).Errorf("GetArtists: Error fetching from dynamo: %s", err)
		return nil, status.Error(codes.Internal, "Error fetching result")
	}

	// check an item was returned
	if data.Item == nil {
		logger.WithField("id", id).Error("GetArtists: No artist found for ID")
		return nil, status.Error(codes.NotFound, "Artist not found")
	}

	// fetch was successful
	logger.WithField("data", data.Item).Info("GetArtists: Artist found")

	// unmarshal response
	res := &setmakerpb.Artist{}
	if err = attributevalue.UnmarshalMap(data.Item, res); err != nil {
		logger.WithField("data", data.Item).Errorf("GetArtists: Could not unmarshal item: %s", err)
		return nil, status.Error(codes.Internal, "Error unmarshaling artist data")
	}

	return res, nil
}

func (d *DynamoRepository) PutArtist(ctx context.Context, artist *setmakerpb.Artist) error {
	// create attribute value map
	item, err := attributevalue.MarshalMap(artist)
	if err != nil {
		logger.WithField("data", artist).Errorf("GetArtists: Could not marshalmap: %s", err)
		return status.Error(codes.InvalidArgument, "Could not map input values for artist")
	}

	// PutItem to dynamo
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TableName),
		Item:      item,
	})
	if err != nil {
		logger.WithField("data", artist).Errorf("GetArtists: Could not PutItem: %s", err)
		return status.Error(codes.Internal, "Failed to persist artist")
	}

	logger.WithField("id", artist.Id).Info("GetArtists: Artist persisted successfully")
	return nil
}

func (d *DynamoRepository) DeleteArtist(ctx context.Context, id uuid.UUID) error {
	logger.WithField("id", id).Infof("GetArtists: Deleting artist")

	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		logger.WithField("id", id).Errorf("GetArtists: Could not delete artist: %s", err)
		return status.Error(codes.Internal, "Artist could not be deleted")
	}

	return nil
}
