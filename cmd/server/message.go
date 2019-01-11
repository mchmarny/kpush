package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mchmarny/kpush/pkg/msg"
	"github.com/mchmarny/kpush/pkg/valid"
)

func signedMessageHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// check method
	if r.Method != http.MethodPost {
		log.Printf("wring method: %s", r.Method)
		http.Error(w, "Invalid method. Only POST supported", http.StatusMethodNotAllowed)
		return
	}

	// parse form to update
	if err := r.ParseForm(); err != nil {
		log.Printf("error parsing form: %v", err)
		http.Error(w, fmt.Sprintf("Post content error (%s)", err),
			http.StatusBadRequest)
		return
	}

	// check for presense of publisher token
	srcToken := r.URL.Query().Get(posterTokenName)
	if srcToken == "" {
		log.Printf("nil token: %s", srcToken)
		http.Error(w, fmt.Sprintf("Invalid request (%s missing)", posterTokenName),
			http.StatusBadRequest)
		return
	}

	// check validity of poster token
	if !contains(knownPublisherTokens, srcToken) {
		log.Printf("invalid token: %s", srcToken)
		http.Error(w, fmt.Sprintf("Invalid publisher token value (%s)", posterTokenName),
			http.StatusBadRequest)
		return
	}

	// parse payload
	payload := &msg.PushData{}
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		log.Printf("error decoding payload: %v", err)
		http.Error(w, fmt.Sprintf("Error decoding payload: %s", err), http.StatusBadRequest)
		return
	}

	// check payload integrity
	if payload.Message == nil || payload.Message.Attributes == nil {
		log.Print("invalid payload")
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// get signature
	payloadSig := payload.Message.Attributes[valid.SignatureAttributeName]
	if payloadSig == "" {
		log.Printf("nil signature: %v", payloadSig)
		http.Error(w, "Invalid payload (missing signature)", http.StatusBadRequest)
		return
	}

	// capture data content
	rawData := payload.Message.Data

	// check signature
	if !valid.IsContentSignatureValid(key, rawData, payloadSig) {
		log.Printf("invalid signature: %s", payloadSig)
		http.Error(w, "Invalid payload signature", http.StatusBadRequest)
		return
	}

	// parse payload data
	pushedMsg, err := msg.MessageFromBytes(rawData)
	if err != nil {
		log.Printf("error parsing payload data: %v", err)
		http.Error(w, fmt.Sprintf("Error decoding payload data: %s", err), http.StatusBadRequest)
		return
	}

	// TODO: do something usefull with the pushed message here
	log.Printf("handler result: %s", pushedMsg)

	// response with the parsed payload data
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(pushedMsg)

}
