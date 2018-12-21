package msg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// SimpleMessage represents a simple message content
type SimpleMessage struct {
	ID        string              `json:"id,omitempty"`
	Timestamp time.Time           `json:"on,omitempty"`
	SourceID  string              `json:"src,omitempty"`
	Data      *SimpleMeasurements `json:"data,omitempty"`
}

// String prints out full message content
func (m SimpleMessage) String() string {
	return fmt.Sprintf("[ID: %s, On: %s, Src: %s, Data: %s]",
		m.ID, m.Timestamp, m.SourceID, m.Data)
}

// Bytes returns byte content of the message
func (m *SimpleMessage) Bytes() []byte {
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(m)
	return reqBodyBytes.Bytes()
}

// SimpleMeasurements represents data readings
type SimpleMeasurements struct {
	Value1 string  `json:"val1,omitempty"`
	Value2 float32 `json:"val2,omitempty"`
	Value3 int     `json:"val3,omitempty"`
}

func (m SimpleMeasurements) String() string {
	return fmt.Sprintf("[Value1: %s, Value2: %f, Value3: %d]",
		m.Value1, m.Value2, m.Value3)
}
