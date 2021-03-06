package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mchmarny/kpush/pkg/msg"
	"github.com/mchmarny/kpush/pkg/valid"
)

func TestPostHandlerUsingSample(t *testing.T) {

	const token = "test-token"
	knownPublisherTokens = []string{token, "other-string-we-are-not-testing"}

	k := []byte("super-long-but-not-so-secure-string")
	key = k

	m := msg.MakeRundomMessage("srcTest")

	b, e := msg.MessageToBytes(m)
	if e != nil {
		t.Errorf("Error message to bytes: %s", e)
		return
	}

	s := valid.MakeSignature(k, b)
	d := &msg.PushData{
		Message: &msg.PushMessage{
			Attributes: make(map[string]string),
			MessageID:  makeID(),
			Data:       b,
		},
	}
	d.Message.Attributes[valid.SignatureAttributeName] = s

	// finally encode the PubSub-like object into bytes
	data, _ := json.Marshal(d)

	req, _ := http.NewRequest("POST", "/push", bytes.NewReader(data))

	q := req.URL.Query()
	q.Add(posterTokenName, token)
	req.URL.RawQuery = q.Encode()

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(signedMessageHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("Handler returned wrong status code: got %v expected %v",
			status, http.StatusAccepted)
		return
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("Error while reading response body %v", err)
		return
	}

	rm := &msg.SimpleMessage{}
	err = json.Unmarshal(body, &rm)
	if err != nil {
		t.Errorf("Error while parsing response body %v", err)
		return
	}

}
