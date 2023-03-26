package akevitt

import (
	"fmt"
)

const accountBucket string = "Accounts"
const gameObjectBucket string = "GameObjects"

type Account struct {
	Username string
	password string
}

func (account Account) Save(key uint64, engine *Akevitt) error {
	return overwriteObject(engine.db, key, accountBucket, account)
}

func (acc Account) Description() string {
	format := "Name: %s\nThis is player.\n"
	return fmt.Sprintf(format, acc.Username)

}
