package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mchmarny/pusheventing/pkg/msg"
	"github.com/mchmarny/pusheventing/cmd/client/pub"
	"github.com/mchmarny/pusheventing/cmd/client/pub/topic"
)

const (
	defaultPushEventingSource = "demoClient"
	defaultPubSubTopicName    = "pusheventing"
	defaultNumberOfMessages   = 3
)

var (
	canceling     bool
	key           []byte
	src           string
	projectID     string
	topicName         string
	numOfMessages int
	publisher     pub.Publisher
)

func main() {

	// flags
	flag.StringVar(&projectID, "project", os.Getenv("GCLOUD_PROJECT"), "Project ID")
	keyStr := flag.String("key", os.Getenv("MSG_SIG_KEY"), "Signature key")
	flag.StringVar(&topicName, "topic", defaultPubSubTopicName, "PubSub topic name [pusheventing]")
	flag.StringVar(&src, "src", defaultPushEventingSource, "Source of data [demoClient]")
	flag.IntVar(&numOfMessages, "messages", defaultNumberOfMessages, "Number of messages to sent [3]")
	flag.Parse()

	if projectID == "" {
		log.Fatal("project flag required (or define GCLOUD_PROJECT env var)")
	}

	if *keyStr == "" {
		log.Fatal("signature key required (or define PUSH_EVENTING_KEY env var)")
	}
	key = []byte(*keyStr)

	// context
	ctx, cancel := context.WithCancel(context.Background())

	// TODO: Make this dynamic based on flag
	p, err := topic.NewPubSubPublisher(ctx, projectID, topicName)
	if err != nil {
		log.Fatalf("error creating publisher: %v", err)
	}

	publisher = p

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
		}else{
			status <- fmt.Sprintf("published msg[%d] %s", i, m.ID)
		}
		sentCount++
	}

	done <- sentCount
	return

}
