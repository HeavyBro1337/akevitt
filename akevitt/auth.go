package akevitt

import "errors"

func (engine *Akevitt) Login(username, password string, session ActiveSession) error {
	account, err := login(username, password, engine.db)
	if err != nil {
		return err
	}
	if isSessionAlreadyActive(*account, &engine.sessions, engine) {
		return errors.New("the session is already active")
	}

	session.SetAccount(account)

	return nil
}

func (engine *Akevitt) Register(username, password string, session ActiveSession) error {
	exists := isAccountExists(username, engine.db)

	if exists {
		return errors.New("account already exists")
	}
	account, err := createAccount(engine.db, username, password)
	session.SetAccount(account)

	return err
}
