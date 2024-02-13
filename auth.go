package akevitt

import (
	"errors"
	"strings"
)

// Login call to the database.
// Note: returns an error if the session is already active.
func (engine *Akevitt) Login(username, password string, session *ActiveSession) error {
	err := validateCredentials(username, password)

	if err != nil {
		return err
	}

	account, err := login(username, password, engine)
	if err != nil {
		return err
	}
	if isSessionAlreadyActive(*account, &engine.sessions, engine) {
		return errors.New("the session is already active")
	}

	session.Account = account

	return nil
}

// Create an account
// Returns an error if account with the same username already exists
func (engine *Akevitt) Register(username, password, repeatPassword string, session *ActiveSession) error {
	err := validateCredentials(username, password)

	if err != nil {
		return err
	}

	if password != repeatPassword {
		return errors.New("passwords don't match")
	}

	exists := isAccountExists(username, engine)

	if exists {
		return errors.New("account already exists")
	}
	account, err := createAccount(engine, username, password)
	session.Account = account

	return err
}

func validateCredentials(username, password string) error {
	username = strings.TrimSpace(username)

	if username == "" {
		return errors.New("username must not be empty")
	}

	if password == "" {
		return errors.New("password must not be empty")
	}

	return nil
}
