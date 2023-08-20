package akevitt

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Hashes password using Bcrypt algorithm
func hashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

// Compares hash and password
func compareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Find[T comparable](collection []T, value T) bool {
	for _, b := range collection {
		if b == value {
			return true
		}
	}
	return false
}

// Converts byte array to `Uint64`
// func byteToInt(source []byte) uint64 {
// 	return binary.BigEndian.Uint64(source)
// }

// Converts `Uint64` to byte array
func intToByte(value uint64) []byte {
	binaryId := make([]byte, 8)
	binary.BigEndian.PutUint64(binaryId, uint64(value))
	return binaryId
}

// Converts `T` to byte array
func serialize[T Object](v T) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// Converts byte array to T struct.
func deserialize[T Object](b []byte) (T, error) {
	var result T
	var decodeBuffer bytes.Buffer
	decodeBuffer.Write(b)
	dec := gob.NewDecoder(&decodeBuffer)
	err := dec.Decode(&result)
	if err != nil {
		return result, err
	}
	return result, err
}

func findByKey[TCollection, T comparable](collection []TCollection, selector func(key TCollection) T, value T) *TCollection {
	if collection == nil {
		panic(errors.New("collection is nil"))
	}
	for _, b := range collection {
		if selector(b) == value {
			return &b
		}
	}
	return nil
}

func RemoveItem[T comparable](l []T, item T) []T {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func MapSlice[T comparable, TResult any](l []T, callback func(v T) TResult) []TResult {
	result := make([]TResult, 0)

	for _, v := range l {
		result = append(result, callback(v))
	}

	return result
}
