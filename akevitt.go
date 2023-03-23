/*
Program written by Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view LICENCE or README
*/

package akevitt

import (
	"akevitt/core/database"
	"akevitt/core/network"
	"errors"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

// The engine struct. Handles connections, user-provided logic and database.
type Akevitt struct {
	activeSessions map[ssh.Session]network.ActiveSession                             // Active sessions
	dbPath         string                                                            // Database path
	rootScreen     func(engine *Akevitt, session network.ActiveSession) *tview.Modal // Screen that user will see on connect
	bind           string                                                            // Port or address to listen
	handle         func(sesh ssh.Session)                                            // Customizable client handle logic
	db             *bolt.DB                                                          // Database file
	mouse          bool                                                              // Allow client to use their mouse
	gameName       string                                                            // Game's title
	callbacks      *engineCallbacks                                                  // Struct for holding all of the callbacks
}
type engineCallbacks struct {
	alreadyLoggedIn func(engine *Akevitt, session network.ActiveSession)
	accountExists   func(engine *Akevitt, session network.ActiveSession)
	loginSuccess    func(engine *Akevitt, session network.ActiveSession)
	loginFail       func(engine *Akevitt, session network.ActiveSession)
	oocMessage      func(engine *Akevitt, session network.ActiveSession, sender network.ActiveSession, message string)
}

// Creates new instance of engine with provided defaults.
func (self *Akevitt) Defaults() *Akevitt {
	self.activeSessions = make(map[ssh.Session]network.ActiveSession)
	self.bind = ":2222"
	self.dbPath = "data/database.db"
	self.gameName = "Change Me!"
	self.callbacks = &engineCallbacks{}
	return self
}

// Listen to specific address and/or port
func (self *Akevitt) Bind(bind string) *Akevitt {
	self.bind = bind

	return self
}

// The file name that will act as a database
func (self *Akevitt) DatabasePath(path string) *Akevitt {
	self.dbPath = path

	return self
}

// Optional: Add an ability to have a mouse interaction.
func (self *Akevitt) UseMouse() *Akevitt {
	self.mouse = true
	return self
}

// Set the game's name. Default: Change Me!
func (self *Akevitt) GameName(name string) *Akevitt {
	self.gameName = name

	return self
}

// If specified, you can provide your own handle logic
func (self *Akevitt) Handle(handle func(sesh ssh.Session)) *Akevitt {
	if handle != nil {
		self.handle = handle
		return self
	}

	self.handle = func(sesh ssh.Session) {
		screen, err := network.NewSessionScreen(sesh)
		if err != nil {
			fmt.Fprintln(sesh.Stderr(), "unable to create screen:", err)
			return
		}

		network.PurgeDeadSessions(&self.activeSessions)

		app := tview.NewApplication().SetScreen(screen).EnableMouse(self.mouse)
		self.activeSessions[sesh] = network.ActiveSession{Chat: nil, Account: nil, UI: app}
		app.SetRoot(self.rootScreen(self, self.activeSessions[sesh]), false)
		if err := app.Run(); err != nil {
			fmt.Fprintln(sesh.Stderr(), err)
			return
		}

		sesh.Exit(0)
	}
	return self
}

// Creates database file if not exists. The custom path must be already specified, before creating.
func (self *Akevitt) CreateDatabaseIfNotExists() *Akevitt {
	db, err := bolt.Open(self.dbPath, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	self.db = db

	return self
}

// Launch the engine
func (self *Akevitt) Run() error {
	defer self.db.Close()

	if self.rootScreen == nil {
		return errors.New("base screen is not provided")
	}
	if self.db == nil {
		return errors.New("database is unused")
	}

	ssh.Handle(self.handle)

	return ssh.ListenAndServe(self.bind, nil)
}

// Set the root screen. The player will see it on connection.
func (self *Akevitt) RootScreen(s func(engine *Akevitt, session network.ActiveSession) *tview.Modal) *Akevitt {
	self.rootScreen = s

	return self
}

func (self *Akevitt) Login(username, password string, session network.ActiveSession) {
	ok, account := database.Login(username, password, self.db)
	if !ok {
		self.callbacks.loginFail(self, session)
		return
	}
	if database.CheckCurrentLogin(*account, &self.activeSessions) {
		self.callbacks.alreadyLoggedIn(self, session)
		return
	}
	session.Account = account
	self.callbacks.loginSuccess(self, session)
	return
}

// Retreives the game's name
func (self *Akevitt) GetGameName() string {
	return self.gameName
}

func (self *Akevitt) ConfigureCallbacks() *engineCallbacks {
	return self.callbacks
}

func (self *engineCallbacks) AlreadyLoggedIn(c func(engine *Akevitt, session network.ActiveSession)) *engineCallbacks {
	self.alreadyLoggedIn = c
	return self
}

func (self *engineCallbacks) LoginSuccesFull(c func(engine *Akevitt, session network.ActiveSession)) *engineCallbacks {
	self.loginSuccess = c
	return self
}

func (self *engineCallbacks) LoginFail(c func(engine *Akevitt, session network.ActiveSession)) *engineCallbacks {
	self.loginSuccess = c
	return self
}

func (self *engineCallbacks) OOCMessage(c func(engine *Akevitt, session network.ActiveSession, sender network.ActiveSession, message string)) *engineCallbacks {
	self.oocMessage = c
	return self
}

func (self *Akevitt) SendOOCMessage(message string, session network.ActiveSession) {
	network.PurgeDeadSessions(&self.activeSessions)
	network.BroadcastMessage(&self.activeSessions, message, session,
		func(message string, sender network.ActiveSession, currentSession network.ActiveSession) {
			self.callbacks.oocMessage(self, session, sender, message)
		})
}
