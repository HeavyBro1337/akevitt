package plugins

import (
	"errors"
	"fmt"
	"strings"

	"github.com/IvanKorchmit/akevitt"
	"golang.org/x/crypto/bcrypt"
)

type AccountPlugin struct {
	engine *akevitt.Akevitt
}

func (plugin *AccountPlugin) login(username string, password string) (*akevitt.Account, error) {
	databasePlugin, err := akevitt.FetchPlugin[akevitt.DatabasePlugin[*akevitt.Account]](plugin.engine)

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

func (plugin *AccountPlugin) isSessionAlreadyActive(acc akevitt.Account, sessions *akevitt.Sessions) bool {
	// We want make sure we purge dead sessions before looking for active.
	akevitt.PurgeDeadSessions(plugin.engine, plugin.engine.GetOnDeadSession())
	for _, v := range *sessions {
		if v.Account == nil {
			continue
		}
		if v.Account.Username == acc.Username {
			return true
		}
	}
	return false
}

// Compares hash and password.
func compareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (plugin *AccountPlugin) Build(engine *akevitt.Akevitt) error {
	plugin.engine = engine
	return nil
}

func (plugin *AccountPlugin) Register(username, password, repeatPassword string, session *akevitt.ActiveSession) error {
	err := validateCredentials(username, password)

	if err != nil {
		return err
	}

	if password != repeatPassword {
		return errors.New("passwords don't match")
	}

	exists := isAccountExists(username, plugin.engine)

	if exists {
		return errors.New("account already exists")
	}
	account, err := createAccount(plugin.engine, username, password)
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

func isAccountExists(username string, engine *akevitt.Akevitt) bool {
	databasePlugin, err := akevitt.FetchPlugin[akevitt.DatabasePlugin[*akevitt.Account]](engine)

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

func createAccount(engine *akevitt.Akevitt, username, password string) (*akevitt.Account, error) {
	exists := isAccountExists(username, engine)

	if exists {
		return nil, fmt.Errorf("account with name %s already exists", username)
	}

	hash, err := hashString(password)

	if err != nil {
		return nil, err
	}

	databasePlugin, err := akevitt.FetchPlugin[akevitt.DatabasePlugin[*akevitt.Account]](engine)

	if err != nil {
		return nil, err
	}

	account := &akevitt.Account{
		Username:       username,
		Password:       hash,
		PersistentData: make(map[string]any),
	}

	err = (*databasePlugin).Save(account)

	return account, err
}

// Hashes password using Bcrypt algorithm
func hashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}
