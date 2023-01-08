package credentials

import (
	"akevitt/core/database/utils"
	"akevitt/core/objects"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"
)

const AccountBucket string = "Accounts"

type Account struct {
	Username string
	Password string
}

// Save, through `gob`, `Account` data at specified key in the database.
func (account Account) Save(key uint64, db *bolt.DB) error {
	errResult := db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(AccountBucket))

		if err != nil {
			return err
		}
		serialized, err := objects.Serialize(account)
		if err != nil {
			return err
		}
		bkt.Put(utils.IntToByte(key), serialized)
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
	return fmt.Sprintf(color.BlueString(format), color.GreenString(acc.Username))

}
