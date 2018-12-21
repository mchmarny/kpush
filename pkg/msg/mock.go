package msg

import (
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

func mockString(length int) string {
	b := make([]byte, length)
	for i := 0; i < length; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

func mockID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("Error while getting id: %v\n", err)
	}
	return id.String()
}

// MakeRundomMessage creates a rundom data loaded message
func MakeRundomMessage(src string) *SimpleMessage {

	return &SimpleMessage{
		ID:        mockID(),
		Timestamp: time.Now(),
		SourceID:  src,
		Data: &SimpleMeasurements{
			Value1: mockString(10),
			Value2: rand.Float32(),
			Value3: rand.Intn(100),
		},
	}

}
