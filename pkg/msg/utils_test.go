package msg

import (
	"testing"
)

func TestMessageBytesUtils(t *testing.T) {

	m1 := MakeRundomMessage("src1")

	// to bytes
	b, err := MessageToBytes(m1)
	if err != nil {
		t.Errorf("Error on converting to bytes: %v", err)
	}

	// from bytes
	m2, err := MessageFromBytes(b)
	if err != nil {
		t.Errorf("Error on converting from bytes: %v", err)
	}

	// validation
	if m1.ID != m2.ID {
		t.Error("Failed to restore ID")
	}

	if m1.Data != m1.Data {
		t.Error("Failed to restore data")
	}

}
