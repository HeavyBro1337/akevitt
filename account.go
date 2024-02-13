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
	databasePlugin, err := FetchPlugin[DatabasePlugin[*Account]](engine)

	if err != nil {
		return err
	}
	return (*databasePlugin).Save(account)
}
