package msg

import (
	"encoding/json"
)

// MessageToBytes returns byte representation of the SimpleMessage
func MessageToBytes(m *SimpleMessage) (b []byte, err error) {
	content, e := json.Marshal(m)
	if e != nil {
		return nil, e
	}
	return content, nil
}

// MessageFromBytes hydrates object from byte content
func MessageFromBytes(content []byte) (m *SimpleMessage, err error) {
	msg := &SimpleMessage{}
	e := json.Unmarshal(content, msg)
	if e != nil {
		return nil, e
	}
	return msg, nil
}
