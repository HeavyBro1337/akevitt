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
	"reflect"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type Akevitt struct {
	sessions      Sessions
	root          UIFunc
	bind          string
	mouse         bool
	dbPath        string
	commands      map[string]CommandFunc
	db            *bolt.DB
	onMessage     MessageFunc
	onDeadSession DeadSessionFunc
	defaultRoom   Room
	rooms         map[uint64]Room
}

// Engine default constructor
func NewEngine() *Akevitt {
	engine := &Akevitt{}
	engine.rooms = make(map[uint64]Room)
	engine.bind = ":2222"
	engine.sessions = make(Sessions)
	engine.dbPath = "data/database.db"
	engine.mouse = false
	engine.commands = make(map[string]CommandFunc)
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

func (engine *Akevitt) ProcessCommand(command string, session ActiveSession) error {
	zeroArg := strings.Fields(command)[0]
	noZeroArgArray := strings.Fields(command)[1:]
	noZeroArg := strings.Join(noZeroArgArray, " ")
	commandFunc, ok := engine.commands[zeroArg]
	if !ok {
		return errors.New("command not found")
	}

	return commandFunc(engine, session, noZeroArg)
}

func (engine *Akevitt) UseMessage(f MessageFunc) *Akevitt {
	engine.onMessage = f

	return engine
}

func (engine *Akevitt) UseSpawnRoom(r Room) *Akevitt {
	engine.defaultRoom = r

	return engine
}

// Provide some callback if session is ended. Note: Some methods are dangerious to call i.e. engine.Message,
// because it may invoke dead session cleanup which will cause stack overflow error and crash the application.
func (engine *Akevitt) UseOnSessionEnd(f DeadSessionFunc) *Akevitt {
	engine.onDeadSession = f
	return engine
}

func saveRoomsRecursively(engine *Akevitt, room Room, visited []string) error {
	if visited == nil {
		visited = make([]string, 0)
	}

	if room == nil {
		return errors.New("room is nil")
	}

	fmt.Printf("Loading Room: %s\n", room.GetName())

	engine.rooms[room.GetKey()] = room

	visited = append(visited, room.GetName())

	for _, v := range room.GetExits() {
		r := v.GetRoom()

		if Find[string](visited, r.GetName()) {
			continue
		}

		err := saveRoomsRecursively(engine, r, visited)

		if err != nil {
			return err
		}
	}
	return nil
}

func (engine *Akevitt) GetCommands() []string {
	result := make([]string, 0)

	for k := range engine.commands {
		result = append(result, k)
	}

	return result
}

func (engine *Akevitt) Lookup(room Room) []GameObject {
	return room.GetObjects()
}

func (engine *Akevitt) GetSpawnRoom() Room {
	return engine.defaultRoom
}

func (engine *Akevitt) GetRoom(key uint64) (Room, error) {
	room, ok := engine.rooms[key]
	if !ok {
		return nil, errors.New("room not found")
	}

	return room, nil
}

func (engine *Akevitt) SaveGameObject(gameObject GameObject, key uint64, account *Account) error {
	return overwriteObject(engine.db, key, account.Username, gameObject)
}

func (engine *Akevitt) SaveObject(gameObject GameObject, key uint64) error {
	return overwriteObject(engine.db, key, "Global", gameObject)
}

func SaveObject[T Object](engine *Akevitt, obj T, category string, key uint64) error {
	return overwriteObject[T](engine.db, key, category, obj)
}

func FindObject[T GameObject](engine *Akevitt, session ActiveSession, key uint64) (T, error) {
	return findObject[T](engine.db, *session.GetAccount(), key)
}

func (engine *Akevitt) Message(channel, message, username string, session ActiveSession) error {
	if engine.onMessage == nil {
		return errors.New("onMessage func is nil")
	}
	purgeDeadSessions(&engine.sessions, engine, engine.onDeadSession)

	for _, v := range engine.sessions {

		err := engine.onMessage(engine, v, channel, message, username)

		if session != v {
			v.GetApplication().Draw()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func Run[TSession ActiveSession](engine *Akevitt) error {
	fmt.Println("Running Akevitt")
	err := createDatabase(engine)

	if err != nil {
		return err
	}

	fmt.Println("Opened database")

	fmt.Println("Loading rooms recursively...")

	err = saveRoomsRecursively(engine, engine.defaultRoom, nil)

	if err != nil {
		return err
	}

	fmt.Println("Done!")

	defer engine.db.Close()

	gob.Register(Account{})

	if engine.root == nil {
		return errors.New("base screen is not provided")
	}

	ssh.Handle(func(sesh ssh.Session) {
		var emptySession TSession
		screen, err := newSessionScreen(sesh)
		if err != nil {
			fmt.Fprintln(sesh.Stderr(), "unable to create screen:", err)
			return
		}
		purgeDeadSessions(&engine.sessions, engine, engine.onDeadSession)
		app := tview.NewApplication().SetScreen(screen).EnableMouse(engine.mouse)

		engine.sessions[sesh] = reflect.New(reflect.TypeOf(emptySession).Elem()).Interface().(TSession)

		engine.sessions[sesh].SetApplication(app)
		engine.sessions[sesh].GetApplication().SetRoot(engine.root(engine, engine.sessions[sesh]), true)
		ticker := time.NewTicker(100 * time.Millisecond)
		quit := make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					purgeDeadSessions(&engine.sessions, engine, engine.onDeadSession)
					app.Draw()
				case <-quit:
					ticker.Stop()
					return
				}
			}
		}()
		if err := app.Run(); err != nil {
			fmt.Fprintln(sesh.Stderr(), err)
			return
		}
		sesh.Exit(0)
	})
	return ssh.ListenAndServe(engine.bind, nil)
}
