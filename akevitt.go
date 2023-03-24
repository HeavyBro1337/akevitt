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
	"reflect"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

// The engine struct. Handles connections, user-provided logic and database.
type Akevitt struct {
	activeSessions map[ssh.Session]network.ActiveSession                                // Active sessions
	dbPath         string                                                               // Database path
	rootScreen     func(engine *Akevitt, session network.ActiveSession) tview.Primitive // Screen that user will see on connect
	bind           string                                                               // Port or address to listen
	handle         func(sesh ssh.Session)                                               // Customizable client handle logic
	db             *bolt.DB                                                             // Database file
	mouse          bool                                                                 // Allow client to use their mouse
	gameName       string                                                               // Game's title
	callbacks      *GameEventHandler                                                    // Struct for holding all of the callbacks
}
type GameEventHandler struct {
	alreadyLoggedIn func(engine *Akevitt, session network.ActiveSession)
	accountExists   func(engine *Akevitt, session network.ActiveSession)
	loginSuccess    func(engine *Akevitt, session network.ActiveSession)
	loginFail       func(engine *Akevitt, session network.ActiveSession)
	oocMessage      func(engine *Akevitt, session network.ActiveSession, sender network.ActiveSession, message string)
	validated       bool // True when it has passed all of the validation
}

// Creates new instance of engine with provided defaults.
func (engine *Akevitt) Defaults() *Akevitt {
	engine.activeSessions = make(map[ssh.Session]network.ActiveSession)
	engine.bind = ":2222"
	engine.dbPath = "data/database.db"
	engine.gameName = "Change Me!"

	return engine
}

// Listen to specific address and/or port
func (engine *Akevitt) Bind(bind string) *Akevitt {
	engine.bind = bind

	return engine
}

// The file name that will act as a database
func (engine *Akevitt) DatabasePath(path string) *Akevitt {
	engine.dbPath = path

	return engine
}

// Optional: Add an ability to have a mouse interaction.
func (engine *Akevitt) UseMouse(toggle bool) *Akevitt {
	engine.mouse = toggle
	return engine
}

// Set the game's name. Default: Change Me!
func (engine *Akevitt) GameName(name string) *Akevitt {
	engine.gameName = name

	return engine
}

// If specified, you can provide your own handle logic
func (engine *Akevitt) Handle(handle func(sesh ssh.Session)) *Akevitt {
	if handle != nil {
		engine.handle = handle
		return engine
	}

	engine.handle = func(sesh ssh.Session) {
		screen, err := network.NewSessionScreen(sesh)
		if err != nil {
			fmt.Fprintln(sesh.Stderr(), "unable to create screen:", err)
			return
		}

		network.PurgeDeadSessions(&engine.activeSessions)

		app := tview.NewApplication().SetScreen(screen).EnableMouse(engine.mouse)
		engine.activeSessions[sesh] = network.ActiveSession{Chat: nil, Account: nil, UI: app}
		engine.activeSessions[sesh].SetRoot(engine.rootScreen(engine, engine.activeSessions[sesh]))
		if err := app.Run(); err != nil {
			fmt.Fprintln(sesh.Stderr(), err)
			return
		}

		sesh.Exit(0)
	}
	return engine
}

// Creates database file if not exists. The custom path must be already specified, before creating.
func (engine *Akevitt) CreateDatabaseIfNotExists() *Akevitt {
	db, err := bolt.Open(engine.dbPath, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	engine.db = db

	return engine
}

// Launch the engine
func (engine *Akevitt) Run() error {
	defer engine.db.Close()

	if engine.rootScreen == nil {
		return errors.New("base screen is not provided")
	}
	if engine.db == nil {
		return errors.New("database is unused")
	}

	ssh.Handle(engine.handle)

	return ssh.ListenAndServe(engine.bind, nil)
}

// Set the root screen. The player will see it on connection.
func (engine *Akevitt) RootScreen(s func(engine *Akevitt, session network.ActiveSession) tview.Primitive) *Akevitt {
	engine.rootScreen = s

	return engine
}

func (engine *Akevitt) Login(username, password string, session network.ActiveSession) {
	ok, account := database.Login(username, password, engine.db)
	if !ok {
		engine.callbacks.loginFail(engine, session)
		return
	}
	if database.CheckCurrentLogin(*account, &engine.activeSessions) {
		engine.callbacks.alreadyLoggedIn(engine, session)
		return
	}
	session.Account = account
	engine.callbacks.loginSuccess(engine, session)
}

// Retreives the game's name
func (engine *Akevitt) GetGameName() string {
	return engine.gameName
}

func (engine *Akevitt) ConfigureCallbacks(event *GameEventHandler) *Akevitt {
	if !event.validated {
		log.Fatal("the event handler is not validated!")
	}

	engine.callbacks = event
	return engine
}

func (event *GameEventHandler) AlreadyLoggedIn(c func(engine *Akevitt, session network.ActiveSession)) *GameEventHandler {
	event.alreadyLoggedIn = c
	return event
}

func (event *GameEventHandler) LoginSuccesFull(c func(engine *Akevitt, session network.ActiveSession)) *GameEventHandler {
	event.loginSuccess = c
	return event
}

func (event *GameEventHandler) LoginFail(c func(engine *Akevitt, session network.ActiveSession)) *GameEventHandler {
	event.loginSuccess = c
	return event
}

func (event *GameEventHandler) OOCMessage(c func(engine *Akevitt, session network.ActiveSession, sender network.ActiveSession, message string)) *GameEventHandler {
	event.oocMessage = c
	return event
}
func (event *GameEventHandler) AccountExists(c func(engine *Akevitt, session network.ActiveSession)) *GameEventHandler {
	event.accountExists = c
	return event
}
func (engine *Akevitt) SendOOCMessage(message string, session network.ActiveSession) {
	network.PurgeDeadSessions(&engine.activeSessions)
	network.BroadcastMessage(&engine.activeSessions, message, session,
		func(message string, sender network.ActiveSession, currentSession network.ActiveSession) {
			engine.callbacks.oocMessage(engine, session, sender, message)
		})
}

func (engine *Akevitt) Register(username, password string, session network.ActiveSession) error {
	if database.DoesAccountExist(username, engine.db) {
		engine.callbacks.accountExists(engine, session)
		return errors.New("account already exists")
	}

	return database.CreateAccount(engine.db, username, password)

}

func (events *GameEventHandler) Finish() {
	// TODO: Implement better validation for detecting errors!
	hasPassed := true
	reflected := reflect.ValueOf(*events)
	fieldNum := reflected.NumField()
	eventType := reflect.TypeOf(*events)
	for i := 0; i < fieldNum; i++ {

		if eventType.Field(i).Type.Kind() != reflect.Func {
			fmt.Println(eventType.Field(i).Type.Kind())
			continue
		}

		if reflected.Field(i).IsNil() {
			log.Printf("error! %s is nil! Did you miss something?\n", eventType.Field(i).Name)
			hasPassed = false
		}

	}
	events.validated = hasPassed
}
