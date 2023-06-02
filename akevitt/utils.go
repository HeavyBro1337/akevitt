package akevitt

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
)

// Converts `Uint64` to byte array
func intToByte(value uint64) []byte {
	binaryId := make([]byte, 8)
	binary.LittleEndian.PutUint64(binaryId, uint64(value))
	return binaryId
}

type Pair[TFirst any, TSecond any] struct {
	First  TFirst
	Second TSecond
}

func filterMap[TKey comparable, TValue any](data map[TKey]TValue, predicate func(k TKey, v TValue) bool) map[TKey]TValue {
	if predicate == nil {
		panic(errors.New("predicate is nil"))
	}

	filtered := make(map[TKey]TValue)

	for k, v := range data {

		if predicate(k, v) {
			filtered[k] = v
		}
	}

	return filtered
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
	return string(result), nil
}

func find[T comparable](collection []T, value T) bool {
	for _, b := range collection {
		if b == value {
			return true
		}
	}
	return false
}
