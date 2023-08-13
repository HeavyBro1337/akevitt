package akevitt

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"

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

// Converts byte array to `Uint64`
// func byteToInt(source []byte) uint64 {
// 	return binary.LittleEndian.Uint64(source)
// }

// Converts `Uint64` to byte array
func intToByte(value uint64) []byte {
	binaryId := make([]byte, 8)
	binary.LittleEndian.PutUint64(binaryId, uint64(value))
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
