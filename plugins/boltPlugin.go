package plugins

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"

	"github.com/IvanKorchmit/akevitt"
	"github.com/boltdb/bolt"
)

type BoltDbPlugin struct {
	path string
	db   *bolt.DB
}

func (plugin *BoltDbPlugin) Build(engine *akevitt.Akevitt) error {
	return createBoltDatabase(plugin)
}

func (plugin *BoltDbPlugin) Save(object akevitt.Object) error {
	return plugin.db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(fmt.Sprint(reflect.TypeOf(object))))

		if err != nil {
			return err
		}

		bytes, err := plugin.serialize(object)

		if err != nil {
			return err
		}

		return bkt.Put([]byte(object.GetName()), bytes)
	})
}

func (plugin *BoltDbPlugin) LoadAll() ([]akevitt.Object, error) {
	objects := make([]akevitt.Object, 0)

	err := plugin.db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(fmt.Sprint(reflect.TypeOf(new(akevitt.Object)))))

		if err != nil {
			return err
		}

		return bkt.ForEach(func(k []byte, v []byte) error {
			var obj akevitt.Object
			err := plugin.deserialize(v, &obj)

			if err != nil {
				return err
			}

			objects = append(objects, obj)

			return nil
		})
	})

	return objects, err
}

func NewBoltPlugin(path string) *BoltDbPlugin {
	return &BoltDbPlugin{
		path: path,
	}
}

func createBoltDatabase(boltPlugin *BoltDbPlugin) error {
	_db, err := bolt.Open(boltPlugin.path, 0600, nil)
	boltPlugin.db = _db

	return err
}

// Converts `T` to byte array.
func (plugin *BoltDbPlugin) serialize(v akevitt.Object) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// Converts byte array to T struct.
func (plugin *BoltDbPlugin) deserialize(b []byte, dest any) error {
	var decodeBuffer bytes.Buffer
	decodeBuffer.Write(b)
	dec := gob.NewDecoder(&decodeBuffer)
	err := dec.Decode(&dest)
	if err != nil {
		return err
	}
	return nil
}
