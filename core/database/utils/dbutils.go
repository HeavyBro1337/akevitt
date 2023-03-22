package utils

import (
	"crypto/sha256"
	"encoding/binary"
)

// Converts `Uint64` to byte array
func IntToByte(value uint64) []byte {
	binaryId := make([]byte, 8)
	binary.LittleEndian.PutUint64(binaryId, uint64(value))
	return binaryId
}

// Hashes string using SHA-256 algorithm
func HashString(input string) string {
	hash := sha256.New()
	_, err := hash.Write([]byte(input))

	if err != nil {
		return ""
	}

	result := hash.Sum(nil)

	return string(result)
}
