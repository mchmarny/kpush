package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mchmarny/kpush/pkg/msg"
	"github.com/mchmarny/kpush/pkg/valid"
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

	// get message byte content
	b, err := msg.MessageToBytes(content)
	if err != nil {
		return fmt.Errorf("Error while converting message to bytes: %v", err)
	}

	// create signature attribute
	attr := make(map[string]string)
	attr[valid.SignatureAttributeName] = valid.MakeSignature(key, b)

	// create a pubsub message with the content attributes
	msg := &msg.PushData{
		Message: &msg.PushMessage{
			MessageID:  content.ID,
			Attributes: attr,
			Data:       b,
		},
	}

	postContent, e := json.Marshal(msg)
	if e != nil {
		return fmt.Errorf("Error while serializing POST payload: %v", err)
	}

	log.Printf("Posting: %s", string(postContent))

	req, err := http.NewRequest("POST", p.url, bytes.NewBuffer(postContent))
	if err != nil {
		return fmt.Errorf("Error on post to %s: %v", p.url, err)
	}

	req.Header.Add("Content-Type", "application/json")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("Error on post exec %s: %v", p.url, err)
	}
	defer resp.Body.Close()

	log.Printf("Response status: %s", resp.Status)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error on body read %s: %v", p.url, err)
	}

	log.Printf("Response body: %s", data)

	return nil

}
