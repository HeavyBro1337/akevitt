package objects

import (
	"bytes"
	"encoding/gob"

	"github.com/boltdb/bolt"
)

// In-game object.
type Object interface {
	Description() string                // Retrieve description about that object
	Save(key uint64, db *bolt.DB) error // Save object into database
}

// Converts `T` to byte array
func Serialize[T Object](v T) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	encodeErr := enc.Encode(v)
	if encodeErr != nil {
		return nil, encodeErr
	}
	return buff.Bytes(), nil
}

// Converts byte array to T struct.
func Deserialize[T Object](b []byte) (T, error) {
	var result T
	var decodeBuffer bytes.Buffer
	decodeBuffer.Write(b)
	dec := gob.NewDecoder(&decodeBuffer)
	err := dec.Decode(&result)
	return result, err
}
