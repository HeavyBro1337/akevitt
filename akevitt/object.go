package akevitt

import (
	"bytes"
	"encoding/gob"
)

type Object interface {
	Description() string                    // Retrieve description about that object
	Save(key uint64, engine *Akevitt) error // Save object into database
}

type GameObject interface {
	Object
	Name() string
	Create(engine *Akevitt, session *ActiveSession, params interface{}) error
	GetMap() map[string]Object
	OnRoomLookup() uint64
	OnLoad(engine *Akevitt) error
}

type Room interface {
	Object
	GetExits() []Exit
	GetKey() uint64
}

type Exit interface {
	Object
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

	return result, err
}

func CreateObject[T GameObject](engine *Akevitt, session *ActiveSession, object T, params interface{}) (T, error) {
	return object, object.Create(engine, session, params)
}
