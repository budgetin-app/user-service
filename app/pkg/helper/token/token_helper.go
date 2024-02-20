package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// SessionTokenLength should be an even number
const SessionTokenLength = 64

func GenerateSessionToken() (string, error) {
	// Determine byte length for the token
	length := SessionTokenLength / 2

	// Generate random bytes
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Hash the random bytes using SHA-256
	hash := sha256.Sum256(randomBytes)

	// Encode the hash to base64 format
	return base64.URLEncoding.EncodeToString(hash[:])[:length], nil
}
