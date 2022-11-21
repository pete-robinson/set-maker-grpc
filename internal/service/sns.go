package service

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	setmakerpb "github.com/pete-robinson/setmaker-proto/dist"
	logger "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type SnsClient struct {
	client   *sns.Client
	topicArn SnsTopic
}

type SnsTopic string

func NewSnsClient(client *sns.Client, topic SnsTopic) *SnsClient {
	return &SnsClient{
		client:   client,
		topicArn: topic,
	}
}

func (s *SnsClient) RaiseArtistCreatedEvent(ctx context.Context, artist *setmakerpb.Artist) error {
	// build the message body
	event := &setmakerpb.Event{
		EventType: setmakerpb.Event_EVENT_ARTIST_CREATED,
		MessageBody: &setmakerpb.Event_ArtistCreated{
			ArtistCreated: &setmakerpb.MessageBody_ArtistCreated{
				Id:   artist.Id,
				Name: artist.Name,
			},
		},
	}

	logger.WithField("MessageBody", event).Infof("Raising event: %s", setmakerpb.Event_EventType_name[int32(event.EventType)])

	// raise the event
	if _, err := s.raise(ctx, event); err != nil {
		return err
	}

	return nil
}

func (s *SnsClient) raise(ctx context.Context, event *setmakerpb.Event) (*string, error) {
	// marshal message
	msg, err := proto.Marshal(event)
	if err != nil {
		logger.WithField("event", event).Errorf("Could not marshal SNS message: %s", err)
		return nil, err
	}

	// build input struct
	snsIn := &sns.PublishInput{
		Message:          aws.String(string(msg)),
		TopicArn:         aws.String(string(s.topicArn)),
	}

	// publish message
	res, err := s.client.Publish(ctx, snsIn)
	if err != nil {
		logger.WithField("event", event).Errorf("Could not publish to SNS topic: %s", err)
		return nil, err
	}

	logger.WithFields(logger.Fields{
		"event":     event,
		"messageId": *res.MessageId,
	}).Info("Published new event to SNS topic")

	return res.MessageId, nil
}
