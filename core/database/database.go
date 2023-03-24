/*
Program written by Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view LICENCE or README
*/

package database

import (
	"akevitt/core/database/utils"
	"akevitt/core/network"
	"akevitt/core/objects"
	"akevitt/core/objects/credentials"
	"errors"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
)

func CreateAccount(db *bolt.DB, username, password string) error {
	var idResult uint64
	account := credentials.Account{Username: username, Password: utils.HashString(password)}

	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return errors.New("invalid data")
	}

	if DoesAccountExist(strings.TrimSpace(account.Username), db) {
		return errors.New("this account already does exist")
	}

	errResult := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(credentials.AccountBucket))
		if err != nil {
			return err
		}
		idResult, _ = bkt.NextSequence()

		serialized, err := objects.Serialize(account)

		if err != nil {
			return err
		}

		bkt.Put(utils.IntToByte(idResult), serialized)
		return nil
	})
	return errResult
}

// Retrieves data, through `gob`, by converting byte array (value) at `key`
// into `Account`.
func GetAccount(key uint64, db *bolt.DB) (account credentials.Account, err error) {
	var result credentials.Account
	dberr := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(credentials.AccountBucket))
		if err != nil {
			return err
		}
		result, err = objects.Deserialize[credentials.Account](bkt.Get(utils.IntToByte(key)))
		if err != nil {
			log.Fatal("Decode error: ", err)
		}

		return nil
	})
	return result, dberr
}

// Checks current account for being in an active sessions. True if the account is already logged in.
func CheckCurrentLogin(acc credentials.Account, sessions *map[ssh.Session]network.ActiveSession) bool {
	// We want make sure we purge dead sessions before looking for active.
	network.PurgeDeadSessions(sessions)
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

// Checks that user exists in the database by username.
func DoesAccountExist(username string, db *bolt.DB) bool {
	var result bool = false
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(credentials.AccountBucket))
		if err != nil {
			return err
		}
		bucket.ForEach(func(k, v []byte) error {
			acc, err := objects.Deserialize[credentials.Account](v)
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
func Login(username string, password string, db *bolt.DB) (bool, *credentials.Account) {
	var accrResult *credentials.Account = nil
	exists := false
	hashedPassword := utils.HashString(password)

	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(credentials.AccountBucket))
		if err != nil {
			return err
		}
		bucket.ForEach(func(k, v []byte) error {
			acc, err := objects.Deserialize[credentials.Account](v)
			if err != nil {
				return err
			}
			if acc.Username == username && acc.Password == hashedPassword {
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
