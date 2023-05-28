package akevitt

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// Converts `Uint64` to byte array
func intToByte(value uint64) []byte {
	binaryId := make([]byte, 8)
	binary.LittleEndian.PutUint64(binaryId, uint64(value))
	return binaryId
}

// Converts byte array to `Uint64`
func byteToInt(source []byte) uint64 {
	return binary.LittleEndian.Uint64(source)
}

// Hashes string using SHA-256 algorithm
func hashString(input string) (string, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(input))

	if err != nil {
		return "", err
	}

	result := hash.Sum(nil)
	fmt.Printf("string(result): %v\n", string(result))
	return string(result), nil
}
