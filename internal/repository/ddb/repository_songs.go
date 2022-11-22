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

type SongList struct {
	Count int32
	Cursor string
	Items []*setmakerpb.Song
}


// Paginated list of songs
func (d *DynamoRepository) ListSongs(ctx context.Context, limit int32, cursor string) (*SongList, error) {
	// decode the cursor
	c, err := utils.DecodeAttributeMap(cursor)
	if err != nil {
		logger.WithField("cursor", cursor).Errorf("ListSongs Repo: Could not decode cursor: %s", err)
		return nil, err
	}

	logger.WithFields(logger.Fields{
		"limit": limit,
		"cursor": cursor,
	}).Info("ListSongs Repo: Scanning dynamo")

	// query DDB
	res, err := d.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(SongsTable),
		Limit: &limit,
		ExclusiveStartKey: c,
	})
	if err != nil {
		logger.Errorf("ListSongs Repo: Error response from dynamo: %s", err)
		return nil, err
	}

	return d.buildPaginatedResponse(res.Items, res.Count, res.LastEvaluatedKey)
}


// Paginated list of songs by artistId
func (d *DynamoRepository) ListSongsByArtist(ctx context.Context, limit int32, cursor string, artistId string) (*SongList, error) {
	// decode the cursor
	c, err := utils.DecodeAttributeMap(cursor)
	if err != nil {
		logger.WithField("cursor", cursor).Errorf("ListSongsByArtist Repo: Could not decode cursor: %s", err)
		return nil, err
	}

	logger.WithFields(logger.Fields{
		"limit": limit,
		"cursor": cursor,
		"artistId": artistId,
	}).Info("ListSongsByArtist Repo: Scanning dynamo")

	// query ddb
	res, err := d.client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(SongsTable),
		IndexName: aws.String("ArtistId-index"),
		KeyConditionExpression: aws.String("ArtistId = :artistId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":artistId": &types.AttributeValueMemberS{
				Value: artistId,
			},
		},
		Limit: &limit,
		ExclusiveStartKey: c,
	})
	if err != nil {
		logger.Errorf("ListSongsByArtist Repo: Error response from dynamo: %s", err)
		return nil, err
	}

	return d.buildPaginatedResponse(res.Items, res.Count, res.LastEvaluatedKey)
}


// Get song by Id
func (d *DynamoRepository) GetSong(ctx context.Context, id uuid.UUID) (*setmakerpb.Song, error) {
	// create key map
	keys, err := attributevalue.MarshalMap(map[string]string{
		"Id": *aws.String(id.String()),
	})
	if err != nil {
		logger.WithField("id", id).Errorf("GetSong Repo: could not marshal map: %s", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid UUID")
	}

	// fetch item from dynamo
	data, err := d.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(SongsTable),
		Key: keys,
	})
	if err != nil {
		logger.WithField("id", id).Errorf("GetSong Repo: Error fetching result from dynamp: %s", err)
		return nil, status.Error(codes.Internal, "error fetching results")
	}

	// check a result was returned
	if data.Item == nil {
		logger.WithField("id", id).Error("GetSong Repo: No song found for ID")
		return nil, status.Error(codes.NotFound, "Song not found")
	}

	logger.WithField("data", data.Item).Info("GetSong Repo: Song found")

	res := &setmakerpb.Song{}
	if err = attributevalue.UnmarshalMap(data.Item, res); err != nil {
		logger.WithField("data", data.Item).Errorf("GetSong Repo: could not unmarshal data: %s", err)
		return nil, status.Error(codes.Internal, "Error unmarshaling song data")
	}

	return res, nil
}


// Put Song
func (d *DynamoRepository) PutSong(ctx context.Context, song *setmakerpb.Song) error {
	// create attribute value map
	item, err := attributevalue.MarshalMap(song)
	if err != nil {
		logger.WithField("song", song).Errorf("PutSong Repo: Could not marshal map: %s", err)
		return status.Error(codes.InvalidArgument, "Could not map input values for song")
	}

	// PutItem to dynamo
	_, err = d.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(SongsTable),
		Item: item,
	})
	if err != nil {
		logger.WithField("song", song).Errorf("PutSong Repo: Could not PutItem: %s", err)
		return status.Error(codes.Internal, "Failed to persist song")
	}

	logger.WithField("id", song.Id).Info("PutSong Repo: Song persisted successfully")
	return nil
}


// Delete song
func (d *DynamoRepository) DeleteSong(ctx context.Context, id uuid.UUID) error {
	logger.WithField("id", id).Infof("DeleteSong Repo: Deleting song")

	_, err := d.client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(SongsTable),
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		logger.WithField("id", id).Errorf("DeleteSong Repo: Could not delete song: %s", err)
		return status.Error(codes.Internal, "Song could not be deleted")
	}

	return nil
}


// build a paginated response of listed songs (by scan or by query)
func (d *DynamoRepository) buildPaginatedResponse(items []map[string]types.AttributeValue, count int32, lek map[string]types.AttributeValue) (*SongList, error) {
	// encode the return cursor
	returnCursor, err := utils.EncodeAttributeMap(lek)
	if err != nil {
		logger.WithField("lastEvaluatedKey", lek).Errorf("buildPaginatedResponse Repo: Unable to encode attribute map: %s", err)
		return nil, err
	}

	// parse results
	var songs []*setmakerpb.Song
	if err = attributevalue.UnmarshalListOfMaps(items, &songs); err != nil {
		logger.Errorf("buildPaginatedResponse Repo: Could not unmarshal results: %s", err)
		return nil, err
	}

	return &SongList{
		Count: count,
		Cursor: string(returnCursor),
		Items: songs,
	}, nil
}
