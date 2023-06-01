/*
Program written by Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view LICENCE or README
*/

package akevitt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
)

func createAccount(db *bolt.DB, username, password string) (*Account, error) {

	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return nil, errors.New("invalid data")
	}
	hashedPassword, err := hashString(password)

	if err != nil {
		return nil, err
	}

	account := Account{Username: username, Password: hashedPassword}

	if doesAccountExist(strings.TrimSpace(account.Username), db) {
		return nil, errors.New("this account already does exist")
	}
	_, err = createObject(db, accountBucket, account)

	return &account, err
}

func overwriteObject[T Object](db *bolt.DB, key uint64, bucket string, object T) error {
	return db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))

		if err != nil {
			return err
		}
		fmt.Printf("object: %v\n", object)
		serialized, err := serialize(object)

		if err != nil {
			return err
		}

		bkt.Put(intToByte(key), serialized)
		return nil
	})
}

func createObject[T Object](db *bolt.DB, bucket string, object T) (uint64, error) {
	var id uint64

	// We obtain new ID and then call overwriteObject
	err := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))

		if err != nil {
			return err
		}

		id, err = bkt.NextSequence()

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}
	return id, overwriteObject(db, id, bucket, object)
}

func getNewKey(db *bolt.DB, bucket string) (uint64, error) {
	var id uint64

	// We obtain new ID and then call overwriteObject
	err := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))

		if err != nil {
			return err
		}

		id, err = bkt.NextSequence()

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Checks current account for being in an active sessions. True if the account is already logged in.
func checkCurrentLogin(acc Account, sessions *map[ssh.Session]*ActiveSession) bool {
	// We want make sure we purge dead sessions before looking for active.
	purgeDeadSession(sessions)
	for _, v := range *sessions {
		if v.Account == nil {
			continue
		}
		if *v.Account == acc {
			return true
		}
	}
	return false
}

func lookupGameObjects(db *bolt.DB, roomKey uint64) ([]GameObject, error) {
	var result []GameObject

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(gameObjectBucket))

		if err != nil {
			return err
		}

		return bucket.ForEach(func(k, v []byte) error {
			obj, err := deserialize[GameObject](v)

			if err != nil {
				return err
			}

			key := obj.OnRoomLookup()

			if key == roomKey {
				result = append(result, obj)
			}

			return nil
		})
	})
	return result, err
}

func findObject[T GameObject](db *bolt.DB, account Account) (T, uint64, error) {
	var id uint64
	var result T
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(gameObjectBucket))
		if err != nil {
			return err
		}
		return bucket.ForEach(func(k, v []byte) error {
			obj, err := deserialize[T](v)
			if err != nil {
				return err
			}
			if account == obj.GetAccount() {
				return nil
			}

			result = obj
			id = byteToInt(k)
			fmt.Printf("result: %v\n", result)
			return nil
		})
	})
	return result, id, err
}

func findObjectByKey[T Object](db *bolt.DB, key uint64, bucket string) (T, error) {
	var result T
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		return bucket.ForEach(func(k, v []byte) error {
			obj, err := deserialize[T](v)
			if err != nil {
				return err
			}
			if byteToInt(k) == key {
				result = obj
				return nil
			}

			result = obj
			return nil
		})
	})
	return result, err
}

// Checks that user exists in the database by username.
func doesAccountExist(username string, db *bolt.DB) bool {
	result := false
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(accountBucket))
		if err != nil {
			return err
		}
		bucket.ForEach(func(k, v []byte) error {
			acc, err := deserialize[*Account](v)
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

func login(username string, password string, db *bolt.DB) (*Account, error) {
	var accref *Account = nil
	hashedPassword, err := hashString(password)

	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(accountBucket))
		if err != nil {
			return err
		}
		bucket.ForEach(func(k, v []byte) error {
			acc, err := deserialize[Account](v)
			if err != nil {
				return err
			}
			fmt.Printf("acc password: %v\n", acc.Password)
			if acc.Username == username && acc.Password == hashedPassword {
				accref = &acc
				return nil
			}
			return nil
		})
		return nil
	})

	if err != nil {
		return nil, err
	}

	if accref == nil {
		return nil, errors.New("wrong name or password")
	}
	return accref, err
}
