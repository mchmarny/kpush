package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	defaultPort             = "8080"
	knownPublisherTokenName = "KNWON_PUBLISHER_TOKENS"
	posterTokenName         = "publisherToken"
)

var (
	knownPublisherTokens []string
	key                  []byte
)

func main() {

	log.Print("Starting server...")

	// handler
	http.HandleFunc("/push", handlePost)

	// token
	tokens := os.Getenv("KNOWN_PUBLISHER_TOKENS")
	if tokens == "" {
		log.Fatalf("KNOWN_PUBLISHER_TOKENS undefined")
	}
	knownPublisherTokens = strings.Split(tokens, ",")

	// key
	keyStr := os.Getenv("MSG_SIG_KEY")
	if keyStr == "" {
		log.Fatalf("MSG_SIG_KEY undefined")
	}
	key = []byte(keyStr)

	// port
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Server listens on port %s \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
