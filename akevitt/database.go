/*
Program written by Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view LICENCE or README
*/

package akevitt

import (
	"errors"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
)

func createAccount(db *bolt.DB, username, password string) (*Account, error) {
	var idResult uint64
	account := Account{Username: username, Password: hashString(password)}

	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return nil, errors.New("invalid data")
	}

	if doesAccountExist(strings.TrimSpace(account.Username), db) {
		return nil, errors.New("this account already does exist")
	}

	errResult := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(accountBucket))
		if err != nil {
			return err
		}
		idResult, _ = bkt.NextSequence()

		serialized, err := serialize(account)

		if err != nil {
			return err
		}

		bkt.Put(intToByte(idResult), serialized)
		return nil
	})
	return &account, errResult
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

// Checks that user exists in the database by username.
func doesAccountExist(username string, db *bolt.DB) bool {
	var result bool = false
	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(accountBucket))
		if err != nil {
			return err
		}
		bucket.ForEach(func(k, v []byte) error {
			acc, err := deserialize[Account](v)
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
func login(username string, password string, db *bolt.DB) (bool, *Account) {
	var accOut *Account = nil
	exists := false
	hashedPassword := hashString(password)

	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(accountBucket))
		if err != nil {
			return err
		}
		bucket.ForEach(func(k, v []byte) error {
			acc, err := deserialize[Account](v)
			if err != nil {
				return err
			}
			if acc.Username == username && acc.Password == hashedPassword {
				accOut = &acc
				exists = true
				return nil
			}
			return nil
		})
		return nil
	})
	return exists, accOut
}
