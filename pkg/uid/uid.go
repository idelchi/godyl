// Package uid provides methods for generating unique identifiers.
package uid

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/google/uuid"
)

// Hash takes a string as input and returns its sha512 hash in hexadecimal format.
func Hash(str string) string {
	hasher := sha512.New()    // Create a new sha512 hasher.
	hasher.Write([]byte(str)) // Write the input string to the hasher.
	hash := hasher.Sum(nil)   // Compute the sha512 hash.

	// Return the hexadecimal encoding of the hash.
	return hex.EncodeToString(hash)
}

// UUID generates a new UUID and returns it as a string.
func UUID() string {
	id := uuid.New()

	return id.String()
}
