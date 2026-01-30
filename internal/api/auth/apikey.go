package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

const keyPrefix = "sk_live_"

// GenerateAPIKey creates a new API key and returns (plainKey, hash, prefix, error)
func GenerateAPIKey() (plain string, hash string, prefix string, err error) {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", "", err
	}

	// Encode as base64 (URL-safe, no padding)
	token := base64.RawURLEncoding.EncodeToString(bytes)
	plain = keyPrefix + token // "sk_live_a1b2c3..."

	// Hash for storage
	hashBytes := sha256.Sum256([]byte(plain))
	hash = hex.EncodeToString(hashBytes[:])

	// Prefix for display (first 12 chars)
	prefix = plain[:12] // "sk_live_a1b2"

	return plain, hash, prefix, nil
}

// HashAPIKey hashes an incoming key for comparison
func HashAPIKey(key string) string {
	hashBytes := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hashBytes[:])
}
