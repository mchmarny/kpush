package pkg

import (
	"strings"
	"testing"
)

func TestMockedMessageSignature(t *testing.T) {

	key := []byte("not-so-secret-test-key")

	msg := mockMessage("src1")
	b1 := msg.Bytes()

	sig1 := MakeSignature(key, b1)

	if !strings.HasPrefix(sig1, signaturePrefix) {
		t.Errorf("Invalid signature format: %s", sig1)
	}

	if !IsContentSignatureValid(key, msg.Bytes(), sig1) {
		t.Error("Invalid message signature")
	}

}
