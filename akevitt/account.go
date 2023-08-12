package akevitt

type Account struct {
	Username string
	Password string
}

func (account *Account) GetKey() uint64 {
	return 0
}

func (account *Account) Save(engine *Akevitt) error {
	return overwriteObject[*Account](engine.db, account.GetKey(), account.Username, account)
}
