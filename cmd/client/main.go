package main

import (
	"log"
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"fmt"
	"cloud.google.com/go/pubsub"

	"github.com/mchmarny/pusheventing/pkg/msg"
	"github.com/mchmarny/pusheventing/pkg/valid"
)

var (
	topic *pubsub.Topic
	key []byte
	canceling  bool
)

func main() {

	// flags
	projectID := flag.String("project", os.Getenv("GCLOUD_PROJECT"), "Project ID")
	topicName := flag.String("topic", "pusheventing", "PubSub topic name [pusheventing]")
	messages := flag.Int("messages", 3, "Number of messages to sent [3]")
	flag.Parse()

	if *projectID == "" {
		log.Fatal("project flag required (or define GCLOUD_PROJECT env var)")
	}

	// context
	ctx, cancel := context.WithCancel(context.Background())

	// pubsub client
	c, err := pubsub.NewClient(ctx, *projectID)
	if err != nil {
		log.Panicf("Failed to create client: %v", err)
	}

	// pubsub topic
	topic = c.Topic(*topicName)

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		log.Println(<-ch)
		canceling = true
		cancel()
		os.Exit(0)
	}()

	results := make(chan string)

	var mu sync.Mutex
	sentCount := 0

	go sendMessages(ctx, *messages)


	for {
		select {
		case <-ctx.Done():
			break
		case r := <-results:
			mu.Lock()
			sentCount++
			fmt.Printf("\a[%d] %s \n", sentCount, r)
			mu.Unlock()
		}
	}



}

func sendMessages(ctx context.Context, n int) string {

	if n < 1 {
		return "Invalid number of messages to send"
	}


	return "M"

}

// Store persist the SimpleEvent
func publish(ctx context.Context, content *msg.SimpleMessage) error {

	// get message byte content
	b := content.Bytes()

	// create signature attribute
	attr := make(map[string]string)
	attr[valid.SignatureAttributeName] = valid.MakeSignature(key, b)

	// create a pubsub message with the content attributes
	msg := &pubsub.Message{
		ID: content.ID,
		Attributes: attr,
		Data: b,
	}

	// publish message
	result := topic.Publish(ctx, msg)

	// check result
	_, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("Error while publishing message: %v", err)
	}

	return nil

}
