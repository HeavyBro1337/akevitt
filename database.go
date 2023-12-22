package akevitt

import (
	"errors"
	"strings"

	"github.com/boltdb/bolt"
)

func createDatabase(engine *Akevitt) error {
	db, err := bolt.Open(engine.dbPath, 0600, nil)
	engine.db = db

	return err
}

func isSessionAlreadyActive(acc Account, sessions *Sessions, engine *Akevitt) bool {
	// We want make sure we purge dead sessions before looking for active.
	purgeDeadSessions(sessions, engine, engine.onDeadSession)
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

func generateKey(db *bolt.DB, bucketName string) (uint64, error) {
	result := uint64(0)

	return result, db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(bucketName))

		if err != nil {
			return err
		}
		key, err := bucket.NextSequence()

		if err != nil {
			return err
		}
		result = key

		return nil
	})
}

func login(username string, password string, db *bolt.DB) (*Account, error) {
	var accref *Account = nil

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(username))
		if bucket == nil {
			return errors.New("wrong name or password")
		}
		acc, err := deserialize[*Account](bucket.Bucket(intToByte(0)).Get(intToByte(0)))
		if err != nil {
			return err
		}
		if compareHash(password, acc.Password) {
			accref = acc
		}
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

func overwriteObject[T Object](db *bolt.DB, key uint64, bucket string, object T) error {
	return db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(bucket))
		if err != nil {
			return err
		}
		serialized, err := serialize(object)
		if err != nil {
			return err
		}
		dataBucket, err := bkt.CreateBucketIfNotExists(intToByte(key))

		if err != nil {
			return err
		}

		return dataBucket.Put(intToByte(0), serialized)
	})
}

func createAccount(db *bolt.DB, username, password string) (*Account, error) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return nil, errors.New("invalid data")
	}
	hashedPassword, err := hashString(password)

	if err != nil {
		return nil, err
	}

	account := &Account{Username: strings.TrimSpace(username), Password: hashedPassword}
	exists := isAccountExists(account.Username, db)

	if exists {
		return nil, errors.New("this account does already exist")
	}
	err = overwriteObject[*Account](db, 0, account.Username, account)
	return account, err
}

func isAccountExists(username string, db *bolt.DB) bool {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(username))
		if bucket != nil {
			return errors.New("account exists")
		}
		return nil
	}) != nil
}

func findObject[T GameObject](db *bolt.DB, account Account, key uint64) (T, error) {
	var result T
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(account.Username))
		if bucket == nil {
			return errors.New("account does not exist")
		}
		dataBucket := bucket.Bucket(intToByte(key))

		if dataBucket == nil {
			return errors.New("object not found")
		}

		r, err := deserialize[T](dataBucket.Get(intToByte(0)))

		result = r
		return err
	})
	return result, err
}
