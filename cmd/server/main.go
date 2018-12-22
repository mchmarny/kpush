package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	defaultPort             = "8080"
	portVariableName        = "PORT"
	knownPublisherTokenName = "KNOWN_PUBLISHER_TOKENS"
	messageSignatureKeyName = "MSG_SIG_KEY"
	posterTokenName         = "publisherToken"
)

var (
	knownPublisherTokens []string
	key                  []byte
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	list := []string{"POST: /post"}
	msg := struct {
		Handlers []string `json:"handlers"`
	}{
		list,
	}

	json.NewEncoder(w).Encode(msg)

}

func main() {

	log.Print("Starting server...")

	// handlers
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/push", handlePost)

	// token
	tokens := os.Getenv(knownPublisherTokenName)
	if tokens == "" {
		log.Fatalf("%s undefined", knownPublisherTokenName)
	}
	knownPublisherTokens = strings.Split(tokens, ",")

	// key
	keyStr := os.Getenv(messageSignatureKeyName)
	if keyStr == "" {
		log.Fatalf("%s undefined", messageSignatureKeyName)
	}
	key = []byte(keyStr)

	// port
	port := os.Getenv(portVariableName)
	if port == "" {
		port = defaultPort
	}

	log.Printf("Server started on port %s \n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
