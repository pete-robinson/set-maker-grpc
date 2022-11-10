package utils

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type AwsConfig struct {
	Region string
}

func BuildAwsConfig(ctx context.Context, config *AwsConfig) (aws.Config, error) {
	cfg := aws.Config{}
	cfg, err := awsconfig.LoadDefaultConfig(ctx, func(opts *awsconfig.LoadOptions) error {
		opts.Region = config.Region
		return nil
	})
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func CreateDynamoClient(cfg aws.Config) *dynamodb.Client {
	client := dynamodb.NewFromConfig(cfg)
	return client
}
