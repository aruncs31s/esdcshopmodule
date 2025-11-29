package persistence

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// UUIDGenerator generates unique IDs.
type UUIDGenerator struct{}

// NewUUIDGenerator creates a new UUIDGenerator.
func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{}
}

// Generate generates a unique ID.
// This is a simplified implementation. In production, use a proper UUID library.
func (g *UUIDGenerator) Generate() string {
	// Generate a simple unique ID based on timestamp and random bytes
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)

	id := make([]byte, 16)
	for i := 0; i < 8; i++ {
		id[i] = byte(timestamp >> (8 * i))
	}
	copy(id[8:], randomBytes)

	return hex.EncodeToString(id)
}
