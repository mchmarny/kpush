package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/mchmarny/kpush/pkg/valid"
)

// MyMessage is a simple struct that represents your message
type MyMessage struct {
	ID        string    `json:"id,omitempty"`
	Timestamp time.Time `json:"on,omitempty"`
	Value     float64   `json:"val,omitempty"`
}

func main() {

	var key string

	fmt.Print("Enter your key: ")
	fmt.Scanf("%s", &key)

	if key == "" {
		log.Fatal("Key required")
	}

	// create key from input
	keyBytes := []byte(key)

	// create instance of your struct
	msg := &MyMessage{
		ID:        "123",
		Timestamp: time.Now(),
		Value:     0.456,
	}

	// convert struct instance to array of bytes
	msgContent, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Error marshaling your object: %v", err)
	}
	log.Printf("Your message is: %s \n", string(msgContent))

	// get signature from the msgContent
	sig := valid.MakeSignature(keyBytes, msgContent)
	log.Printf("Your message signature is: %s \n", sig)

	// validate that signature is valid
	if !valid.IsContentSignatureValid(keyBytes, msgContent, sig) {
		log.Fatal("Invalid message signature")
	}

	log.Print("Done")

}
