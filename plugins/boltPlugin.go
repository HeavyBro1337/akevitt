package plugins

import (
	"bytes"
	"encoding/gob"

	"github.com/IvanKorchmit/akevitt"
	"github.com/boltdb/bolt"
)

var db *bolt.DB = nil

type BoltDbPlugin[T akevitt.Object] struct {
	path string
}

func (plugin *BoltDbPlugin[T]) Build(engine *akevitt.Akevitt) error {
	return createBoltDatabase(plugin)
}

func (plugin *BoltDbPlugin[T]) Save(object T) error {
	return db.Update(func(tx *bolt.Tx) error {
		bytes, err := plugin.serialize(object)

		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(bytes)

		return err
	})
}

func (plugin *BoltDbPlugin[T]) LoadAll() ([]T, error) {
	objects := make([]T, 0)

	err := db.Update(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			obj, err := plugin.deserialize(name)

			if err != nil {
				return err
			}

			objects = append(objects, obj)

			return nil
		})
	})

	return objects, err
}

func NewBoltPlugin[T akevitt.Object](path string) *BoltDbPlugin[T] {
	return &BoltDbPlugin[T]{
		path: path,
	}
}

func createBoltDatabase[T akevitt.Object](boltPlugin *BoltDbPlugin[T]) error {
	_db, err := bolt.Open(boltPlugin.path, 0600, nil)
	db = _db

	return err
}

// Converts `T` to byte array.
func (plugin *BoltDbPlugin[T]) serialize(v T) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// Converts byte array to T struct.
func (plugin *BoltDbPlugin[T]) deserialize(b []byte) (T, error) {
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
