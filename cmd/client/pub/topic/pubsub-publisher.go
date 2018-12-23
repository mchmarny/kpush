package topic

import (
	"context"
	"errors"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"

	"github.com/mchmarny/kpush/pkg/msg"
	"github.com/mchmarny/kpush/pkg/valid"
)

// PubSubPublisher represents GCP PubSub publisher
type PubSubPublisher struct {
	topic *pubsub.Topic
}

// NewPubSubPublisher creates a configured version of the GCP PubSub publisher
func NewPubSubPublisher(ctx context.Context, projectID, topicName string) (p *PubSubPublisher, err error) {

	if projectID == "" || topicName == "" {
		return nil, errors.New("Invalid project or topic arguments")
	}

	// pubsub client
	c, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		log.Panicf("Failed to create client: %v", err)
	}

	publisher := &PubSubPublisher{
		topic: c.Topic(topicName),
	}

	return publisher, nil

}

// Publish publishes messages to GCP PubSub
func (p *PubSubPublisher) Publish(ctx context.Context, key []byte, content *msg.SimpleMessage) error {

	// get message byte content
	b, err := msg.MessageToBytes(content)
	if err != nil {
		return fmt.Errorf("Error while converting message to bytes: %v", err)
	}

	// create signature attribute
	attr := make(map[string]string)
	attr[valid.SignatureAttributeName] = valid.MakeSignature(key, b)

	// create a pubsub message with the content attributes
	msg := &pubsub.Message{
		ID:         content.ID,
		Attributes: attr,
		Data:       b,
	}

	// publish message
	result := p.topic.Publish(ctx, msg)

	// check result
	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("Error while publishing message: %v", err)
	}

	return nil

}
