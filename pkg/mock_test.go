package pkg

import (
	"testing"
)

func TestMockedMessageUniqueness(t *testing.T) {

	msg1 := mockMessage("src1")
	msg2 := mockMessage("src2")

	if msg1.ID == msg2.ID {
		t.Error("Failed to generate unique IDs")
	}

	if msg1.SourceID == msg2.SourceID {
		t.Error("Failed to generate unique src ID")
	}

	if msg1.Data.Value1 == msg2.Data.Value1 {
		t.Error("Failed to generate unique val1")
	}

	if msg1.Data.Value2 == msg2.Data.Value2 {
		t.Error("Failed to generate unique val2")
	}

	if msg1.Data.Value3 == msg2.Data.Value3 {
		t.Error("Failed to generate unique val3")
	}

}
