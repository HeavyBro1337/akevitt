/*
Program written by Maxwell Jensen, Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view the man page or README.md
*/

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"github.com/boltdb/bolt"
)

const BUCKET_ACCOUNTS string = "Accounts"

type Account struct {
	Username string
	Password string
}

// Save, through `gob`, `Account` data at specified key in the database.
func (account Account) save(key uint64, db *bolt.DB) error {
	errResult := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(BUCKET_ACCOUNTS))

		if err != nil {
			return err
		}
		serialized, err := account.Serialize()
		if err != nil {
			return err
		}
		bkt.Put(intToByte(key), serialized)
		return nil
	})
	return errResult
}

func (account Account) String() string {
	return fmt.Sprintf("%s %s", account.Username, account.Password)
}

func createAccount(db *bolt.DB, account Account) error {
	errResult := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(BUCKET_ACCOUNTS))
		if err != nil {
			return err
		}
		id, _ := bkt.NextSequence()

		serialized, err := account.Serialize()

		if err != nil {
			return err
		}

		bkt.Put(intToByte(id), serialized)
		return nil
	})
	return errResult
}

// Retrieves data, through `gob`, by converting byte array (value) at `key`
// into `Account`.
func getAccount(key uint64, db *bolt.DB) Account {
	var result Account
	db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(BUCKET_ACCOUNTS))
		if err != nil {
			return err
		}
		var decodeBuffer bytes.Buffer
		decodeBuffer.Write(bkt.Get(intToByte(key)))
		dec := gob.NewDecoder(&decodeBuffer)
		decErr := dec.Decode(&result)
		if decErr != nil {
			log.Fatal("Decode error: ", decErr)
		}
		return nil
	})
	return result
}

// Converts `Uint64` to byte array
func intToByte(value uint64) []byte {
	binaryId := make([]byte, 8)
	binary.LittleEndian.PutUint64(binaryId, uint64(value))
	return binaryId
}

// Converts `Account` to byte array
func (account Account) Serialize() ([]byte, error) {
	var accBuff bytes.Buffer
	enc := gob.NewEncoder(&accBuff)
	encodeErr := enc.Encode(account)
	if encodeErr != nil {
		return nil, encodeErr
	}
	return accBuff.Bytes(), nil
}
