package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mchmarny/pusheventing/pkg/msg"
	"github.com/mchmarny/pusheventing/pkg/valid"
)

func TestPostHandler(t *testing.T) {

	token := "test-token"
	knownPublisherTokens = []string{token, "other-string-we-are-not-testing"}

	k := []byte("not-so-secret-test-key")
	key = k

	m := msg.MakeRundomMessage("srcTest")
	b := m.Bytes()
	s := valid.MakeSignature(k, b)
	d := &msg.PushData{
		Message: &msg.PushMessage{
			Attributes: make(map[string]string),
			MessageID:  makeID(),
			Data:       []byte(base64.StdEncoding.EncodeToString(b)),
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
	handler := http.HandlerFunc(handlePost)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("Handler returned wrong status code: got %v expected %v",
			status, http.StatusAccepted)
	}

	body, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Errorf("Error while reading response body %v", err)
	}

	rm := &msg.SimpleMessage{}
	err = json.Unmarshal(body, &rm)
	if err != nil {
		t.Errorf("Error while parsing response body %v", err)
	}

	if rm.ID != m.ID {
		t.Errorf("Invalid Message ID: got: %s expected %s", rm.ID, m.ID)
	}

}
