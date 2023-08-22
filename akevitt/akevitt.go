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
	onDialogue    DialogueFunc
	defaultRoom   Room
	rooms         map[uint64]Room
}

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

func (engine *Akevitt) GetCommands() []string {
	result := make([]string, 0)

	for k := range engine.commands {
		result = append(result, k)
	}

	return result
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
