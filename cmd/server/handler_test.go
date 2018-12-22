package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mchmarny/pusheventing/pkg/msg"
)

func TestPostHandlerUsingSample(t *testing.T) {

	const testFilePath = "../../samples/push-payload.json"
	const token = "test-token"
	knownPublisherTokens = []string{token, "other-string-we-are-not-testing"}

	data, err := getFileContent(testFilePath)
	if err != nil {
		t.Errorf("Error while opening %s: %v", testFilePath, err)
		return
	}

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
