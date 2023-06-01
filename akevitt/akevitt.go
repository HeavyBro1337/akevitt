/*
Program written by Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view LICENCE or README
*/

package akevitt

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

// The engine struct. Handles connections, user-provided logic and database.
type Akevitt struct {
	activeSessions map[ssh.Session]*ActiveSession                                                 // Active sessions
	dbPath         string                                                                         // Database path
	rootScreen     func(engine *Akevitt, session *ActiveSession) tview.Primitive                  // Screen that user will see on connect
	bind           string                                                                         // Port or address to listen
	db             *bolt.DB                                                                       // Database file
	mouse          bool                                                                           // Allow client to use their mouse
	gameName       string                                                                         // Game's title
	hooks          *GameEventHandler                                                              // Struct for holding all of the callbacks
	defaultRoom    Room                                                                           // Default room where new players will spawn.
	commands       map[string]func(engine *Akevitt, session *ActiveSession, command string) error // Registered commands
}
type GameEventHandler struct {
	oocMessage  func(engine *Akevitt, session *ActiveSession, sender *ActiveSession, message string)
	roomMessage func(engine *Akevitt, session *ActiveSession, sender *ActiveSession, message string)
	validated   bool // True when it has passed all of the validation
}

// Creates new instance of engine with provided defaults.
func (engine *Akevitt) UseDefaults() *Akevitt {
	engine.activeSessions = make(map[ssh.Session]*ActiveSession)
	engine.bind = ":2222"
	engine.dbPath = "data/database.db"
	engine.gameName = "Change Me!"
	engine.commands = make(map[string]func(engine *Akevitt, session *ActiveSession, command string) error)
	return engine
}

// Listen to specific address and/or port
func (engine *Akevitt) UseBindAddress(bind string) *Akevitt {
	engine.bind = bind

	return engine
}

// The file name that will act as a database
func (engine *Akevitt) UseDatabasePath(path string) *Akevitt {
	engine.dbPath = path

	return engine
}

// Optional: Add an ability to have a mouse interaction.
func (engine *Akevitt) UseMouse(toggle bool) *Akevitt {
	engine.mouse = toggle
	return engine
}

// Set the game's name. Default: Change Me!
func (engine *Akevitt) UseGameName(name string) *Akevitt {
	engine.gameName = name

	return engine
}

// Creates database file if not exists. The custom path must be already specified, before creating.
func (engine *Akevitt) UseCreateDatabaseIfNotExists() *Akevitt {
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

	ssh.Handle(func(sesh ssh.Session) {
		screen, err := newSessionScreen(sesh)
		if err != nil {
			fmt.Fprintln(sesh.Stderr(), "unable to create screen:", err)
			return
		}

		purgeDeadSession(&engine.activeSessions)

		app := tview.NewApplication().SetScreen(screen).EnableMouse(engine.mouse)
		engine.activeSessions[sesh] = &ActiveSession{Chat: nil, Account: nil, UI: app}

		engine.activeSessions[sesh].RelatedGameObjects = make(map[string]GameObject)

		engine.activeSessions[sesh].SetRoot(engine.rootScreen(engine, engine.activeSessions[sesh]))
		if err := app.Run(); err != nil {
			fmt.Fprintln(sesh.Stderr(), err)
			return
		}

		sesh.Exit(0)
	})

	return ssh.ListenAndServe(engine.bind, nil)
}

// Set the root screen. The player will see it on connection.
func (engine *Akevitt) UseRootScreen(s func(engine *Akevitt, session *ActiveSession) tview.Primitive) *Akevitt {
	engine.rootScreen = s

	return engine
}

func (engine *Akevitt) Login(username, password string, session *ActiveSession) error {
	account, err := login(username, password, engine.db)
	if err != nil {
		return err
	}
	if checkCurrentLogin(*account, &engine.activeSessions) {
		return errors.New("the session is already active")
	}
	session.Account = account
	return nil
}

// Retreives the game's name
func (engine *Akevitt) GetGameName() string {
	return engine.gameName
}

func (engine *Akevitt) ConfigureCallbacks(event *GameEventHandler) *Akevitt {
	if !event.validated {
		log.Fatal("the event handler is not validated!")
	}

	engine.hooks = event
	return engine
}

func (event *GameEventHandler) OOCMessage(c func(engine *Akevitt, session *ActiveSession, sender *ActiveSession, message string)) *GameEventHandler {
	event.oocMessage = c
	return event
}

func (event *GameEventHandler) Message(c func(engine *Akevitt, session *ActiveSession, sender *ActiveSession, message string)) *GameEventHandler {
	event.roomMessage = c
	return event
}

func (engine *Akevitt) SendOOCMessage(message string, session *ActiveSession) {
	purgeDeadSession(&engine.activeSessions)
	broadcastMessage(engine.activeSessions, message, session,
		func(message string, sender *ActiveSession, currentSession *ActiveSession) {
			engine.hooks.oocMessage(engine, currentSession, sender, message)
		})
}

func (engine *Akevitt) WhisperMessage(message string, session *ActiveSession, receiver *ActiveSession, onMessage func(message string, sender *ActiveSession, currentSession *ActiveSession)) {
	onMessage(message, session, receiver)
}

func (engine *Akevitt) SendRoomMessage(message string, session *ActiveSession) {
	purgeDeadSession(&engine.activeSessions)
	broadcastMessage(engine.activeSessions, message, session,
		func(message string, sender *ActiveSession, currentSession *ActiveSession) {
			engine.hooks.roomMessage(engine, currentSession, sender, message)
		})
}

func (engine *Akevitt) Register(username, password string, session *ActiveSession) error {
	if doesAccountExist(username, engine.db) {

		return errors.New("account already exists")
	}
	account, err := createAccount(engine.db, username, password)
	session.Account = account
	return err
}

func FindObject[T GameObject](engine *Akevitt, session *ActiveSession) (T, uint64, error) {
	return findObject[T](engine.db, *session.Account)
}

func (engine *Akevitt) RegisterCommand(command string, function func(e *Akevitt, session *ActiveSession, command string) error) *Akevitt {
	command = strings.TrimSpace(command)
	engine.commands[command] = function

	return engine
}

func (engine *Akevitt) ProcessCommand(command string, session *ActiveSession) error {
	zeroarg := strings.Fields(command)[0]
	nozeroargarr := strings.Fields(command)[1:]
	nozeroarg := strings.Join(nozeroargarr, " ")
	commandFunc, ok := engine.commands[zeroarg]

	if !ok {
		return errors.New("command not found")
	}

	commandFunc(engine, session, nozeroarg)

	return nil
}

func (events *GameEventHandler) Finish() {
	// TODO: Implement better validation for detecting errors!
	hasPassed := true
	reflected := reflect.ValueOf(*events)
	fieldNum := reflected.NumField()
	eventType := reflect.TypeOf(*events)
	for i := 0; i < fieldNum; i++ {

		if eventType.Field(i).Type.Kind() != reflect.Func {
			continue
		}

		if reflected.Field(i).IsNil() {
			log.Printf("error! %s is nil! Did you miss something?\n", eventType.Field(i).Name)
			hasPassed = false
		}

	}
	events.validated = hasPassed
}

func (engine *Akevitt) SaveObject(gameObject GameObject, key uint64) error {
	return overwriteObject(engine.db, key, gameObjectBucket, gameObject)
}

func (engine *Akevitt) SaveWorldObject(object Object, key uint64) error {
	return overwriteObject(engine.db, key, worldObjectsBucket, object)
}

func (engine *Akevitt) SetSpawnRoom(room Room) *Akevitt {
	engine.defaultRoom = room

	key, err := getNewKey(engine.db, worldObjectsBucket)

	if err != nil {
		log.Fatal(err)
	}

	engine.defaultRoom.Save(key, engine)

	return engine
}

func (engine *Akevitt) GetSpawnRoom() Room {
	return engine.defaultRoom
}

func (engine *Akevitt) GetNewKey(isWorld bool) (uint64, error) {
	if !isWorld {
		return getNewKey(engine.db, gameObjectBucket)
	} else {
		return getNewKey(engine.db, worldObjectsBucket)
	}
}

func (engine *Akevitt) GetObject(key uint64) (Object, error) {
	return findObjectByKey[GameObject](engine.db, key)
}

func (engine *Akevitt) Lookup(roomKey uint64) ([]GameObject, error) {
	return lookupGameObjects(engine.db, roomKey)
}
