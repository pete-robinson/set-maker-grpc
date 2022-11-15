package service

import "github.com/aws/aws-sdk-go-v2/service/sns"

type SnsClient struct {
	client *sns.Client
}

func NewSnsClient(client *sns.Client) *SnsClient {
	return &SnsClient{
		client: client,
	}
}
