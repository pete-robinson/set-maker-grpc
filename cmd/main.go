package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	repository "github.com/pete-robinson/set-maker-grpc/internal/repository/ddb"
	"github.com/pete-robinson/set-maker-grpc/internal/service"
	transport "github.com/pete-robinson/set-maker-grpc/internal/transport/grpc"
	"github.com/pete-robinson/set-maker-grpc/internal/utils"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	EnvAwsAccessKey    = "AWS_ACCESS_KEY_ID"
	EnvAwsAccessSecret = "AWS_SECRET_ACCESS_KEY"
	EnvAwsRegion       = "AWS_REGION"
	EnvSnsTopic        = "EVENT_TOPIC"
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

	// init repository
	dynamoClient := utils.CreateDynamoClient(awsConfig)
	repo := repository.NewDynamoRepository(dynamoClient)

	// sns service
	snsClient := utils.CreateSnsClient(awsConfig)
	t := service.SnsTopic(os.Getenv(EnvSnsTopic))
	sns := service.NewSnsClient(snsClient, t)

	// init Service
	svc := service.NewService(repo, sns)

	// init GRPC Server
	server, err := transport.NewServer(svc)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	setmakerpb.RegisterSetMakerServiceServer(s, server)
	reflection.Register(s)

	err = utils.RunGrpcServer(ctx, s)
	if err != nil {
		panic(err)
	}

}
