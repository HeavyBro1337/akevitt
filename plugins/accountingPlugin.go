package plugins

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/IvanKorchmit/akevitt/engine"
	"github.com/rivo/tview"
	"golang.org/x/crypto/bcrypt"
)

type AccountPlugin struct {
	engine *engine.Akevitt
}

func NewAccountPlugin() *AccountPlugin {
	return &AccountPlugin{}
}

func (plugin *AccountPlugin) login(username, password string) (*engine.Account, error) {
	databasePlugin, err := engine.FetchPlugin[engine.DatabasePlugin[*engine.Account]](plugin.engine)

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

func (plugin *AccountPlugin) LoginSession(username, password string, session *engine.ActiveSession) error {
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

func (plugin *AccountPlugin) isSessionAlreadyActive(acc engine.Account, sessions *engine.Sessions) bool {
	// We want make sure we purge dead sessions before looking for active.
	engine.PurgeDeadSessions(plugin.engine, plugin.engine.GetOnDeadSession()...)
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

func (plugin *AccountPlugin) Build(engine *engine.Akevitt) error {
	plugin.engine = engine
	return nil
}

func (plugin *AccountPlugin) Register(username, password, repeatPassword string, session *engine.ActiveSession) error {
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

func isAccountExists(username string, eng *engine.Akevitt) bool {
	databasePlugin, err := engine.FetchPlugin[engine.DatabasePlugin[*engine.Account]](eng)

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

func createAccount(eng *engine.Akevitt, username, password string) (*engine.Account, error) {
	exists := isAccountExists(username, eng)

	if exists {
		return nil, fmt.Errorf("account with name %s already exists", username)
	}

	hash, err := hashString(password)

	if err != nil {
		return nil, err
	}

	databasePlugin, err := engine.FetchPlugin[engine.DatabasePlugin[*engine.Account]](eng)

	if err != nil {
		return nil, err
	}

	account := &engine.Account{
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

func RegistrationScreen(eng *engine.Akevitt, session *engine.ActiveSession, nextScreen engine.UIFunc) tview.Primitive {
	username := ""
	password := ""
	repeatPassword := ""

	account := engine.FetchPluginUnsafe[*AccountPlugin](eng)

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
				engine.ErrorBox(err.Error(), session.Application, form)
				return
			}

			session.Application.SetRoot(nextScreen(eng, session), true)
		})

	return form
}

func LoginScreen(eng *engine.Akevitt, session *engine.ActiveSession, nextScreen engine.UIFunc) tview.Primitive {
	username := ""
	password := ""

	form := tview.NewForm()

	account := engine.FetchPluginUnsafe[*AccountPlugin](eng)

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
				engine.ErrorBox(err.Error(), session.Application, form)
				return
			}

			session.Application.SetRoot(nextScreen(eng, session), true)
		})

	form.SetTitle("Login")

	return form
}