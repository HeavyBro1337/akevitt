package plugins

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/IvanKorchmit/akevitt"
	"github.com/rivo/tview"
	"golang.org/x/crypto/bcrypt"
)

type AccountPlugin struct {
	engine *akevitt.Akevitt
}

func (plugin *AccountPlugin) login(username, password string) (*akevitt.Account, error) {
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

func (plugin *AccountPlugin) LoginSession(username, password string, session *akevitt.ActiveSession) error {
	acc, err := plugin.login(username, password)
	if err != nil {
		return err
	}

	sessions := plugin.engine.GetSessions()

	if plugin.isSessionAlreadyActive(*acc, &sessions) {
		return errors.New("the session is already active")
	}
	session.Account = acc
	return nil
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

func RegistrationScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession, nextScreen akevitt.UIFunc) tview.Primitive {
	username := ""
	password := ""
	repeatPassword := ""

	account := akevitt.FetchPluginUnsafe[*AccountPlugin](engine)

	form := tview.NewForm()

	form.AddInputField("Username", "", 0, func(textToCheck string, lastChar rune) bool {
		if !unicode.IsLetter(lastChar) && !unicode.IsDigit(lastChar) || lastChar > unicode.MaxASCII {
			return false
		}

		username = textToCheck
		return true
	}, nil).
		AddPasswordField("Repeat password", "", 0, '*', func(text string) {
			password = text
		}).
		AddPasswordField("Repeat password", "", 0, '*', func(text string) {
			repeatPassword = text
		}).
		AddButton("Register", func() {
			err := account.Register(username, password, repeatPassword, session)

			if err != nil {
				akevitt.ErrorBox(err.Error(), session.Application, form)
				return
			}

			session.Application.SetRoot(nextScreen(engine, session), true)
		})

	return form
}

func LoginScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession, nextScreen akevitt.UIFunc) tview.Primitive {
	username := ""
	password := ""

	form := tview.NewForm()

	account := akevitt.FetchPluginUnsafe[*AccountPlugin](engine)

	form.AddInputField("Username", "", 0, func(textToCheck string, lastChar rune) bool {
		if !unicode.IsLetter(lastChar) && !unicode.IsDigit(lastChar) || lastChar > unicode.MaxASCII {
			return false
		}

		username = textToCheck
		return true
	}, nil).
		AddPasswordField("Password", "", 0, '*', func(text string) {
			password = text
		}).
		AddButton("Login", func() {
			err := account.LoginSession(username, password, session)

			if err != nil {
				akevitt.ErrorBox(err.Error(), session.Application, form)
				return
			}

			session.Application.SetRoot(nextScreen(engine, session), true)
		})

	form.SetTitle("Login")

	return form
}
