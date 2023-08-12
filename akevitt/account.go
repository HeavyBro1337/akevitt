package akevitt

type Account struct {
	Username string
	Password string
}

func (account *Account) Save(engine *Akevitt) error {
	return overwriteObject[*Account](engine.db, 0, account.Username, account)
}
