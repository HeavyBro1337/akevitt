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
	"github.com/gliderlabs/ssh"
)

const BUCKET_ACCOUNTS string = "Accounts"

type Account struct {
	Username string
	Password string
}

// Save, through `gob`, `Account` data at specified key in the database.
func (account Account) Save(key uint64, db *bolt.DB) error {
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
	// Do not ever pass the password.
	return fmt.Sprintf("%s", account.Username)
}

func createAccount(db *bolt.DB, account Account) (id uint64, err error) {
	var idResult uint64
	if doesAccountExists(account.Username, db) {
		return 0, nil
	}
	errResult := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(BUCKET_ACCOUNTS))
		if err != nil {
			return err
		}
		idResult, _ = bkt.NextSequence()

		serialized, err := account.Serialize()

		if err != nil {
			return err
		}
		
		bkt.Put(intToByte(idResult), serialized)
		return nil
	})
	return idResult, errResult
}

// Retrieves data, through `gob`, by converting byte array (value) at `key`
// into `Account`.
func getAccount(key uint64, db *bolt.DB) (account Account, err error) {
	var result Account
	dberr := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(BUCKET_ACCOUNTS))
		if err != nil {
			return err
		}
		result, err = DeserializeAccount(bkt.Get(intToByte(key)))
		if err != nil {
			log.Fatal("Decode error: ", err)
		}
		
		return nil
	})
	return result, dberr
}
// Checks current account for being in an active sessions. True if the account is already logged in.
func checkCurrentLogin(acc Account, sessions *map[ssh.Session]*Account) bool {
	// We want make sure we purge dead sessions before looking for active. 
	purgeDeadSessions(sessions)
	for _, v := range *sessions {
		if v == nil {
			continue
		}
		if *v == acc {
			return true
		}
	}
	return false
}
// Checks that user exists in the database by username.
func doesAccountExists(username string, db *bolt.DB) bool {
	var result bool = false
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKET_ACCOUNTS))
		if err != nil {
			return err
		}
		bucket.ForEach(func(k, v []byte) error {
			acc, err := DeserializeAccount(v)
			if err != nil {
				return err
			}
			if acc.Username == username {
				result = true
				return nil
			}
			return nil
		})
		return nil
	})
	return result
}


// Logins character and retrieves account from database. It returns true if the login was successfull
func Login(username string, password string, db *bolt.DB) (bool, *Account)  {
	var accrResult *Account = nil
	exists := false
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKET_ACCOUNTS))
		if err != nil {
			return err
		}
		bucket.ForEach(func(k, v []byte) error {
			acc, err := DeserializeAccount(v)
			if err != nil {
				return err
			}
			if acc.Username == username && acc.Password == password {
				accrResult = &acc
				exists = true
				return nil
			}
			return nil
		})
		return nil
	})
	return exists, accrResult
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
// Converts byte array to Account struct.
func DeserializeAccount(b []byte) (Account, error) {
	var result Account
	var decodeBuffer bytes.Buffer
	decodeBuffer.Write(b)
	dec := gob.NewDecoder(&decodeBuffer)
	err := dec.Decode(&result)
	return result, err
}