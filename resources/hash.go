package resources

import (
	"crypto/sha256"
	"encoding/hex"
)

// Returns a hashed SHA256 string from the given string.
func Hash(b []byte) (string, error) {
	h := sha256.New()
	_, err := h.Write(b)
	if err != nil {
		return "", err
	}

	r := h.Sum(nil)
	s := hex.EncodeToString(r)

	return s, nil
}
