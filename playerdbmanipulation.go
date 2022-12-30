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

func (p Player) String() string {
	return fmt.Sprintf("%s %s", p.Username, p.Password)
}

func CreatePlayer(db *bolt.DB, player Player) error {
	errResult := db.Update(func(tx *bolt.Tx) error {
		b, berr := tx.CreateBucketIfNotExists([]byte(BUCKET_USERS))
		if berr != nil {
			return berr
		}
		var plrBuff bytes.Buffer
		enc := gob.NewEncoder(&plrBuff) // Will write to network.
		encodeErr := enc.Encode(player)
		if encodeErr != nil {
			return encodeErr
		}
		id, _ := b.NextSequence()
		binaryId := make([]byte, 8)
		binary.LittleEndian.PutUint64(binaryId, uint64(id))
		b.Put(binaryId, plrBuff.Bytes())
		return nil
	})
	return errResult
}
func GetPlayer(id uint64, db *bolt.DB) Player {
	var result Player
	db.Update(func(tx *bolt.Tx) error {
		b, berr := tx.CreateBucketIfNotExists([]byte(BUCKET_USERS))
		if berr != nil {
			return berr
		}
		var decodeBuffer bytes.Buffer
		binaryId := make([]byte, 8)
		binary.LittleEndian.PutUint64(binaryId, uint64(id))
		decodeBuffer.Write(b.Get(binaryId))
		dec := gob.NewDecoder(&decodeBuffer)
		decErr := dec.Decode(&result)
		if decErr != nil {
			log.Fatal("Decode error: ", decErr)
		}
		return nil
	})
	return result
}
