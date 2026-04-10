package opaquestr

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate creates a Base64-URL encoded string using a
// random byte slice of some specified size.
func Generate(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	str := base64.RawURLEncoding.EncodeToString(b)
	return str, nil
}
