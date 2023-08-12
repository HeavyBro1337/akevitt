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

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type akevitt struct {
	sessions Sessions
	root     UIFunc
	bind     string
	mouse    bool
	dbPath   string
	db       *bolt.DB
}

// Engine default constructor
func NewEngine() *akevitt {
	engine := &akevitt{}
	engine.sessions = make(Sessions)
	engine.dbPath = "data/database.db"
	engine.mouse = false
	return engine
}

func (engine *akevitt) UseBind(bindAddress string) *akevitt {
	engine.bind = bindAddress

	return engine
}

func (engine *akevitt) UseMouse() *akevitt {
	engine.mouse = true

	return engine
}

func (engine *akevitt) Run() error {
	gob.Register(Account{})
	defer engine.db.Close()
	if engine.root == nil {
		return errors.New("base screen is not provided")
	}
	if engine.db == nil {
		return errors.New("database is unused")
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
