package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mchmarny/kpush/cmd/client/pub"
	"github.com/mchmarny/kpush/cmd/client/pub/http"
	"github.com/mchmarny/kpush/cmd/client/pub/topic"
	"github.com/mchmarny/kpush/pkg/msg"
)

const (
	defaultPushEventingSource = "demoClient"
	defaultPubSubTopicName    = "kpush"
	defaultNumberOfMessages   = 3
)

var (
	canceling     bool
	key           []byte
	src           string
	url           string
	projectID     string
	topicName     string
	numOfMessages int
	publisher     pub.Publisher
)

func main() {

	// flags
	flag.StringVar(&projectID, "project", os.Getenv("GCLOUD_PROJECT"), "Project ID")
	keyStr := flag.String("key", os.Getenv("MSG_SIG_KEY"), "Signature key")
	flag.StringVar(&topicName, "topic", defaultPubSubTopicName, "PubSub topic name [kpush]")
	flag.StringVar(&src, "src", defaultPushEventingSource, "Source of data [demoClient]")
	flag.StringVar(&url, "url", "", "(Optional) Target URL where data will be sent directly, no topic")
	flag.IntVar(&numOfMessages, "messages", defaultNumberOfMessages, "Number of messages to sent [3]")
	flag.Parse()

	if projectID == "" {
		log.Fatal("project flag required (or define GCLOUD_PROJECT env var)")
	}

	if *keyStr == "" {
		log.Fatal("signature key required (or define MSG_SIG_KEY env var)")
	}
	key = []byte(*keyStr)

	// context
	ctx, cancel := context.WithCancel(context.Background())

	var err error
	if url == "" {
		log.Printf("Using topic publisher: %s:%s", projectID, topicName)
		publisher, err = topic.NewPubSubPublisher(ctx, projectID, topicName)
		if err != nil {
			log.Fatalf("error creating PubSub publisher: %v", err)
		}
	} else {
		log.Printf("Using http publisher: %s", url)
		publisher, err = http.NewHTTPPublisher(ctx, url)
		if err != nil {
			log.Fatalf("error creating HTTP publisher: %v", err)
		}
	}

	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		log.Println(<-ch)
		canceling = true
		cancel()
		os.Exit(0)
	}()

	status := make(chan string)
	done := make(chan int)

	// start sending data
	go sendMessages(ctx, status, done)

F:
	for {
		select {
		case <-ctx.Done():
			break F
		case s := <-status:
			fmt.Printf("   status: %s \n", s)
		case c := <-done:
			fmt.Printf("sent %d messages \n", c)
			break F
		}
	}

}

func sendMessages(ctx context.Context, status chan<- string, done chan<- int) {

	// loop through to start sending messages
	sentCount := 0
	for i := 0; i < numOfMessages; i++ {
		m := msg.MakeRundomMessage(src)
		err := publisher.Publish(ctx, key, m)
		if err != nil {
			status <- fmt.Sprintf("error msg[%d] %s", i, err.Error())
		} else {
			status <- fmt.Sprintf("published msg[%d] %s", i, m.ID)
		}
		sentCount++
	}

	done <- sentCount
	return

}
