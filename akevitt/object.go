package akevitt

import (
	"bytes"
	"encoding/gob"
)

type Object interface {
	Description() string                    // Retrieve description about that object
	Save(key uint64, engine *Akevitt) error // Save object into database
}

// Gameobject which is associated with the account.
type GameObject interface {
	Object
	Create(engine *Akevitt, session *ActiveSession, params interface{}) error
	GetAccount() Account
	GetRoom() Room
}

type Room interface {
	Object
	GetExits() []Exit
	GetKeys(engine *Akevitt) ([]uint64, error)
	AddObject(engine *Akevitt, key uint64) error
	RemoveObject(engine *Akevitt, key uint64) error
}

type Exit interface {
	Object
}

// Converts `T` to byte array
func serialize[T Object](v T) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	encodeErr := enc.Encode(v)
	if encodeErr != nil {
		return nil, encodeErr
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
	return result, err
}

func CreateObject[T GameObject](engine *Akevitt, session *ActiveSession, object T, params interface{}) (T, error) {
	return object, object.Create(engine, session, params)
}