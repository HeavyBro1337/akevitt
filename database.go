package akevitt

import (
	"errors"
	"fmt"
)

func isSessionAlreadyActive(acc Account, sessions *Sessions, engine *Akevitt) bool {
	// We want make sure we purge dead sessions before looking for active.
	PurgeDeadSessions(engine, engine.onDeadSession)
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

func login(username string, password string, engine *Akevitt) (*Account, error) {
	databasePlugin, err := FetchPlugin[DatabasePlugin[*Account]](engine)

	if err != nil {
		return nil, err
	}

	accounts, err := (*databasePlugin).LoadAll()

	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		if account.Username != username {
			continue
		}

		if !compareHash(password, account.Password) {
			continue
		}

		return account, nil
	}

	return nil, errors.New("wrong username or password")
}

func createAccount(engine *Akevitt, username, password string) (*Account, error) {
	exists := isAccountExists(username, engine)

	if exists {
		return nil, fmt.Errorf("account with name %s already exists", username)
	}

	hash, err := hashString(password)

	if err != nil {
		return nil, err
	}

	databasePlugin, err := FetchPlugin[DatabasePlugin[*Account]](engine)

	if err != nil {
		return nil, err
	}

	account := &Account{
		Username: username,
		Password: hash,
	}

	err = (*databasePlugin).Save(account)

	return account, err
}

func isAccountExists(username string, engine *Akevitt) bool {
	databasePlugin, err := FetchPlugin[DatabasePlugin[*Account]](engine)

	if err != nil {
		panic(err)
	}

	accounts, err := (*databasePlugin).LoadAll()

	if err != nil {
		panic(err)
	}

	for _, account := range accounts {
		if account.Username == username {
			return true
		}
	}

	return false
}
