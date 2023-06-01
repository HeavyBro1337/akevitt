package akevitt

import (
	"fmt"
)

const accountBucket string = "Accounts"
const gameObjectBucket string = "GameObjects"
const worldObjectsBucket string = "WorldObjects"

type Account struct {
	Username string
	Password string
}

func (account Account) Save(key uint64, engine *Akevitt) error {
	return overwriteObject(engine.db, key, accountBucket, account)
}

func (account Account) Name() string {
	return account.Username
}

func (account Account) Description() string {
	format := "Name: %s\nThis is player.\n"
	return fmt.Sprintf(format, account.Username)

}
