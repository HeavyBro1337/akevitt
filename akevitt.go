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
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

// The engine instance which can be passed as an argument and provide some useful methods like
// Login, Register, Message, Dialogue, etc.
// Methods with name starting like Use should be called in a main function during the initialisation step.
// To actually run the engine, you must call Run function and pass the engine instance
// Example: fmt.Fatal(akevitt.Run[*MySessionStruct](engine))
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
	onDialogue    DialogueFunc
	defaultRoom   Room
	rooms         map[uint64]Room
	heartbeats    map[int]*pair[time.Ticker, []func() error]
}

// Execute the command specified in a `command`.
// The command can be registered using the useRegisterCommand method.
// Returns an error if the given command not found or the result of associated function returns an error.
func (engine *Akevitt) ExecuteCommand(command string, session ActiveSession) error {
	zeroArg := strings.Fields(command)[0]
	noZeroArgArray := strings.Fields(command)[1:]
	noZeroArg := strings.Join(noZeroArgArray, " ")
	commandFunc, ok := engine.commands[zeroArg]
	if !ok {
		return errors.New("command not found")
	}

	return commandFunc(engine, session, noZeroArg)
}

// Gets currently registered commands. This is useful if your game implements auto-completion.
func (engine *Akevitt) GetCommands() []string {
	result := make([]string, 0)

	for k := range engine.commands {
		result = append(result, k)
	}

	return result
}

// Get sspawn room if specified. Useful for setting character's initial room during its creation.
func (engine *Akevitt) GetSpawnRoom() Room {
	return engine.defaultRoom
}

// Obtains currently loaded rooms by key. It will return an error if room not found.
func (engine *Akevitt) GetRoom(key uint64) (Room, error) {
	room, ok := engine.rooms[key]
	if !ok {
		return nil, errors.New("room not found")
	}

	return room, nil
}

func (engine *Akevitt) startHeartBeats(interval int) {
	go func() {
		t, ok := engine.heartbeats[interval]
		errResults := make([]int, 0)
		if !ok {
			LogWarn(fmt.Sprintf("ticker %d does not exist", interval))
			return
		}
		for range t.f.C {
			for i, v := range t.s {
				if v == nil {
					continue
				}
				if v() != nil {
					errResults = append(errResults, i)
				}
			}

			for i := len(errResults) - 1; i >= 0; i-- {
				t.s = RemoveItemByIndex(t.s, i)
			}
		}

	}()

}

// Run the given instance of engine.
// You should pass your own implementation of ActiveSession,
// so it can be controlled of how your game would behave
func (engine *Akevitt) Run() error {
	fmt.Println("Running Akevitt")
	err := createDatabase(engine)
	if err != nil {
		return err
	}

	fmt.Println("Opened database")

	fmt.Println("Loading rooms recursively...")

	err = saveRoomsRecursively(engine, engine.defaultRoom, nil)

	for k := range engine.heartbeats {
		engine.startHeartBeats(k)
	}

	if err != nil {
		return err
	}

	fmt.Println("Done!")

	defer engine.db.Close()

	gob.Register(Account{})

	if engine.root == nil {
		return errors.New("base screen is not provided")
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				purgeDeadSessions(&engine.sessions, engine, engine.onDeadSession)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	ssh.Handle(func(sesh ssh.Session) {
		emptySession := ActiveSession{}
		screen, err := newSessionScreen(sesh)
		if err != nil {
			fmt.Fprintln(sesh.Stderr(), "unable to create screen:", err)
			return
		}
		purgeDeadSessions(&engine.sessions, engine, engine.onDeadSession)
		app := tview.NewApplication().SetScreen(screen).EnableMouse(engine.mouse)

		emptySession.Application = app

		engine.sessions[sesh] = emptySession

		emptySession.Application.SetRoot(engine.root(engine, engine.sessions[sesh]), true)

		go func() {
			for {
				select {
				case <-ticker.C:
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
	usePubKey := ssh.HostKeyFile("id_rsa")

	allowKeys := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return true
	})

	return ssh.ListenAndServe(engine.bind, nil, allowKeys, usePubKey)
}
