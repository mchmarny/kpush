package msg

import (
	"fmt"
	"time"
)

// PushData represents the raw data submitted to handler
type PushData struct {
	Message *PushMessage `json:"message"`
}

// PushMessage represents the message portion of PushData
type PushMessage struct {
	Attributes map[string]string `json:"attributes"`
	MessageID  string            `json:"messageId"`
	Data       []byte            `json:"data"`
}

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
