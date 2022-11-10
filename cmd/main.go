package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/pete-robinson/set-maker-grpc/internal/grpc/api"
	repository "github.com/pete-robinson/set-maker-grpc/internal/repository/ddb"
	"github.com/pete-robinson/set-maker-grpc/internal/service"
	transport "github.com/pete-robinson/set-maker-grpc/internal/transport/grpc"
	"github.com/pete-robinson/set-maker-grpc/internal/utils"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	EnvAwsAccessKey    = "aws_access_key_id"
	EnvAwsAccessSecret = "aws_secret_access_key"
	EnvAwsRegion       = "aws_region"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Errorf("BOOT ERROR. COULD NOT LOAD CONFIG: %s", err)
		panic(err)
	}

	// init context
	ctx := context.Background()

	// build AWS config obj
	awsConfigObj := &utils.AwsConfig{
		Region: os.Getenv(EnvAwsRegion),
	}

	awsConfig, err := utils.BuildAwsConfig(ctx, awsConfigObj)
	if err != nil {
		logger.Errorf("BOOT ERROR. COULD NOT BUILD AWS CONFIG: %s", err)
		panic(err)
	}

	dynamoClient := utils.CreateDynamoClient(awsConfig)

	// init epository
	repo := repository.NewDynamoRepository(dynamoClient)

	// init Service
	service := service.NewService(repo)

	// init GRPC Server
	server, err := transport.NewServer(service)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	api.RegisterSetMakerServiceServer(grpcServer, server)

	err = utils.RunGrpcServer(ctx, grpcServer)
	if err != nil {
		panic(err)
	}

}
