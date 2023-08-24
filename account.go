package akevitt

// Basic structure for storing credential information.
// After registering an account, the password is hashed in a proper way
// To create one you would need to invoke `engine.Register(username, password, session)`
type Account struct {
	Username string
	Password string
}

// Save account into a database
func (account *Account) Save(engine *Akevitt) error {
	return overwriteObject[*Account](engine.db, 0, account.Username, account)
}
