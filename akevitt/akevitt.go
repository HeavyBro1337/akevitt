/*
Program written by Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view LICENCE or README
*/

package akevitt

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type Akevitt struct {
	sessions  Sessions
	root      UIFunc
	bind      string
	mouse     bool
	dbPath    string
	commands  map[string]CommandFunc
	db        *bolt.DB
	onMessage MessageFunc
}

// Engine default constructor
func NewEngine() *Akevitt {
	engine := &Akevitt{}
	engine.bind = ":2222"
	engine.sessions = make(Sessions)
	engine.dbPath = "data/database.db"
	engine.mouse = false
	return engine
}

func (engine *Akevitt) UseBind(bindAddress string) *Akevitt {
	engine.bind = bindAddress

	return engine
}

func (engine *Akevitt) UseRootUI(uiFunc UIFunc) *Akevitt {
	engine.root = uiFunc

	return engine
}

func (engine *Akevitt) UseDBPath(path string) *Akevitt {
	engine.dbPath = path

	return engine
}

func (engine *Akevitt) UseMouse() *Akevitt {
	engine.mouse = true

	return engine
}

func (engine *Akevitt) RegisterCommand(command string, function CommandFunc) *Akevitt {
	command = strings.TrimSpace(command)
	engine.commands[command] = function
	return engine
}

func (engine *Akevitt) Login(username, password string, session *ActiveSession) error {
	account, err := login(username, password, engine.db)
	if err != nil {
		return err
	}
	if checkCurrentLogin(*account, &engine.sessions) {
		return errors.New("the session is already active")
	}
	session.Account = account
	return nil
}

func (engine *Akevitt) Register(username, password string, session *ActiveSession) error {
	exists := isAccountExists(username, engine.db)

	if exists {
		return errors.New("account already exists")
	}
	account, err := createAccount(engine.db, username, password)
	session.Account = account
	return err
}

func (engine *Akevitt) ProcessCommand(command string, session *ActiveSession) error {
	zeroArg := strings.Fields(command)[0]
	noZeroArgArray := strings.Fields(command)[1:]
	noZeroArg := strings.Join(noZeroArgArray, " ")
	commandFunc, ok := engine.commands[zeroArg]
	if !ok {
		return errors.New("command not found")
	}

	return commandFunc(engine, session, noZeroArg)
}

func (engine *Akevitt) UseOnMessage(f MessageFunc) *Akevitt {
	engine.onMessage = f

	return engine
}

func (engine *Akevitt) Run() error {
	fmt.Println("Running Akevitt")

	err := createDatabase(engine)

	if err != nil {
		log.Fatal(err)
	}

	defer engine.db.Close()

	gob.Register(Account{})

	if engine.root == nil {
		return errors.New("base screen is not provided")
	}

	ssh.Handle(func(sesh ssh.Session) {
		screen, err := newSessionScreen(sesh)
		if err != nil {
			fmt.Fprintln(sesh.Stderr(), "unable to create screen:", err)
			return
		}
		purgeDeadSessions(&engine.sessions)
		app := tview.NewApplication().SetScreen(screen).EnableMouse(engine.mouse)
		engine.sessions[sesh] = &ActiveSession{Account: nil, UI: app}
		engine.sessions[sesh].UI.SetRoot(engine.root(engine, engine.sessions[sesh]), true)
		if err := app.Run(); err != nil {
			fmt.Fprintln(sesh.Stderr(), err)
			return
		}
		sesh.Exit(0)
	})
	return ssh.ListenAndServe(engine.bind, nil)
}
