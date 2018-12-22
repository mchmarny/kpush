package valid

import (
	"strings"
	"testing"

	"github.com/mchmarny/pusheventing/pkg/msg"
)

func TestMockedMessageSignature(t *testing.T) {

	key := []byte("not-so-secret-test-key")

	m := msg.MakeRundomMessage("src1")
	b1, err := msg.MessageToBytes(m)
	if err != nil {
		t.Errorf("Error getting bytes from message: %v", err)
	}

	sig1 := MakeSignature(key, b1)

	if !strings.HasPrefix(sig1, signaturePrefix) {
		t.Errorf("Invalid signature format: %s", sig1)
	}

	if !IsContentSignatureValid(key, b1, sig1) {
		t.Error("Invalid message signature")
	}

}
