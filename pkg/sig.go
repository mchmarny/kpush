package pkg

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"log"
	"strings"
)

const (
	signaturePrefix = "sha1="
	signatureLength = 45 // len(SignaturePrefix) + len(hex(sha1))
)

// MakeSignature makes content signature from content payload
func MakeSignature(key, content []byte) string {
	dst := make([]byte, 40)
	computed := hmac.New(sha1.New, key)
	computed.Write(content)
	hex.Encode(dst, computed.Sum(nil))
	return signaturePrefix + string(dst)
}

func computedContentSignature(key, content []byte) []byte {
	computed := hmac.New(sha1.New, key)
	computed.Write(content)
	return []byte(computed.Sum(nil))
}

// IsContentSignatureValid computes new content signature and compares it to the original one
func IsContentSignatureValid(key []byte, content []byte, signature string) bool {

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		log.Printf("Invalid signature format: %s", signature)
		return false
	}

	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(computedContentSignature(key, content), actual)
}
