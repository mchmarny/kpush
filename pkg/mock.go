package pkg

import (
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/mchmarny/pusheventing/pkg/types"
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

func mockMessage(src string) *types.SimpleMessage {

	return &types.SimpleMessage{
		ID:        mockID(),
		Timestamp: time.Now(),
		SourceID:  src,
		Data: &types.SimpleMeasurements{
			Value1: mockString(10),
			Value2: rand.Float32(),
			Value3: rand.Intn(100),
		},
	}

}
