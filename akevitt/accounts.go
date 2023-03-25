package akevitt

import (
	"fmt"

	"github.com/boltdb/bolt"
)

const accountBucket string = "Accounts"

type Account struct {
	Username string
	Password string
}

// Save, through `gob`, `Account` data at specified key in the database.
func (account Account) Save(key uint64, db *bolt.DB) error {
	errResult := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(accountBucket))

		if err != nil {
			return err
		}
		serialized, err := serialize(account)
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
	return account.Username
}

func (acc Account) Description() string {
	format := "Name: %s\nThis is player.\n"
	return fmt.Sprintf(format, acc.Username)

}
