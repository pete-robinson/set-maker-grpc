package repository

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	ArtistsTable = "artists"
	SongsTable   = "songs"
)

type DynamoRepository struct {
	client *dynamodb.Client
}

func NewDynamoRepository(client *dynamodb.Client) *DynamoRepository {
	return &DynamoRepository{
		client: client,
	}
}
