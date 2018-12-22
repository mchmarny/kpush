package topic

import (
	"context"
	"errors"

	"github.com/mchmarny/pusheventing/pkg/msg"
)

// HTTPPublisher basic HTTP publisher
type HTTPPublisher struct {
	url string
}

// NewHTTPPublisher creates a configured version of the HTTP publisher
func NewHTTPPublisher(ctx context.Context, url string) (p *HTTPPublisher, err error) {

	if url == "" {
		return nil, errors.New("Invalid url arguments")
	}

	publisher := &HTTPPublisher{
		url: url,
	}

	return publisher, nil

}

// Publish publishes messages to HTTP URL
func (p *HTTPPublisher) Publish(ctx context.Context, key []byte, content *msg.SimpleMessage) error {

	return errors.New("Not yet implemented")

}
