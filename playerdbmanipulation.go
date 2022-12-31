package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const BUCKET_USERS string = "Players"

type Player struct {
	Username string
	Password string
}

func (p Player) Save(id uint64, db *bolt.DB) error {
	errResult := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BUCKET_USERS))

		if err != nil {
			return err
		}
		serialized, err := p.Serialize()
		if err != nil {
			return err
		}
		b.Put(Int2Byte(id), serialized)
		return nil
	})
	return errResult
}
func (p Player) String() string {
	return fmt.Sprintf("%s %s", p.Username, p.Password)
}

func CreatePlayer(db *bolt.DB, player Player) error {
	errResult := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BUCKET_USERS))
		if err != nil {
			return err
		}
		id, _ := b.NextSequence()

		serialized, err := player.Serialize()

		if err != nil {
			return err
		}

		b.Put(Int2Byte(id), serialized)
		return nil
	})
	return errResult
}
func GetPlayer(id uint64, db *bolt.DB) Player {
	var result Player
	db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(BUCKET_USERS))
		if err != nil {
			return err
		}
		var decodeBuffer bytes.Buffer
		decodeBuffer.Write(b.Get(Int2Byte(id)))
		dec := gob.NewDecoder(&decodeBuffer)
		decErr := dec.Decode(&result)
		if decErr != nil {
			log.Fatal("Decode error: ", decErr)
		}
		return nil
	})
	return result
}

// Converts Uint64 to byte array
func Int2Byte(value uint64) []byte {
	binaryId := make([]byte, 8)
	binary.LittleEndian.PutUint64(binaryId, uint64(value))
	return binaryId
}

func (p Player) Serialize() ([]byte, error) {
	var plrBuff bytes.Buffer
	enc := gob.NewEncoder(&plrBuff)
	encodeErr := enc.Encode(p)
	if encodeErr != nil {
		return nil, encodeErr
	}
	return plrBuff.Bytes(), nil
}
