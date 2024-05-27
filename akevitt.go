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

	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type Akevitt struct {
	sessions      Sessions
	root          UIFunc
	bind          string
	mouse         bool
	dbPath        string
	initFunc      []func(*ActiveSession)
	commands      map[string]CommandFunc
	onDeadSession DeadSessionFunc
	onDialogue    DialogueFunc
	defaultRoom   *Room
	rooms         map[uint64]*Room
	plugins       []Plugin
	rsaKey        string
	heartbeats    map[int]*Pair[time.Ticker, []func() error]
}

// Execute the command specified in a `command`.
// The command can be registered using the useRegisterCommand method.
// Returns an error if the given command not found or the result of associated function returns an error.
func (engine *Akevitt) ExecuteCommand(command string, session *ActiveSession) error {
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
func (engine *Akevitt) GetSpawnRoom() *Room {
	return engine.defaultRoom
}

// Obtains currently loaded rooms by key. It will return an error if room not found.
func (engine *Akevitt) GetRoom(key uint64) (*Room, error) {
	room, ok := engine.rooms[key]
	if !ok {
		return nil, errors.New("room not found")
	}

	return room, nil
}

func (engine *Akevitt) GetSessions() Sessions {
	return engine.sessions
}

func (engine *Akevitt) GetOnDeadSession() DeadSessionFunc {
	return engine.onDeadSession
}

// Run the given instance of engine.
// You should pass your own implementation of ActiveSession,
// so it can be controlled of how your game would behave
func (engine *Akevitt) Run() error {
	fmt.Println("Running Akevitt")

	fmt.Println("Building plugins...")

	for _, plugin := range engine.plugins {
		if err := plugin.Build(engine); err != nil {
			fmt.Println("Build failed...")
			return err
		}
	}

	fmt.Println("Done!")

	fmt.Println("Loading rooms recursively...")

	err := saveRoomsRecursively(engine, engine.defaultRoom, nil)

	if err != nil {
		return err
	}

	fmt.Println("Done!")

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
				PurgeDeadSessions(engine, engine.onDeadSession)
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
		PurgeDeadSessions(engine, engine.onDeadSession)
		app := tview.NewApplication().SetScreen(screen).EnableMouse(engine.mouse)

		emptySession.Application = app
		emptySession.Data = make(map[string]any)

		if engine.initFunc != nil {
			for _, fn := range engine.initFunc {
				fn(&emptySession)
			}
		}

		engine.sessions[sesh] = &emptySession

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
	usePubKey := ssh.HostKeyFile(engine.rsaKey)

	allowKeys := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		return true
	})

	return ssh.ListenAndServe(engine.bind, nil, allowKeys, usePubKey)
}

func saveRoomsRecursively(engine *Akevitt, room *Room, visited []string) error {
	if visited == nil {
		visited = make([]string, 0)
	}

	if room == nil {
		return errors.New("room is nil")
	}

	fmt.Printf("Loading Room: %s\n", room.Name)

	engine.rooms[room.GetKey()] = room

	visited = append(visited, room.Name)

	for _, v := range room.Exits {
		r := v.Room

		if Find[string](visited, r.Name) {
			continue
		}

		err := saveRoomsRecursively(engine, r, visited)

		if err != nil {
			return err
		}
	}
	return nil
}
