package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/google/uuid"
)

func getFileContent(path string) (conten []byte, err error) {

	jf, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jf.Close()
	data, _ := ioutil.ReadAll(jf)
	return data, nil
}

func contains(list []string, val string) bool {
	for _, item := range list {
		if item == val {
			return true
		}
	}
	return false
}

func makeID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Error while getting id: %v\n", err)
	}
	return id.String()
}

func mustGetEnv(key, fallbackValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	if fallbackValue == "" {
		log.Fatalf("Required env var (%s) not set", key)
	}

	return fallbackValue
}
