package akevitt

import "strings"

// Specify an address to listen.
// Example: :22, 127.0.0.1:2222, etc.
func (engine *Akevitt) UseBind(bindAddress string) *Akevitt {
	engine.bind = bindAddress

	return engine
}

// Accepts function which returns the UI root screen.
func (engine *Akevitt) UseRootUI(uiFunc UIFunc) *Akevitt {
	engine.root = uiFunc

	return engine
}

// Specify path to save database
func (engine *Akevitt) UseDBPath(path string) *Akevitt {
	engine.dbPath = path

	return engine
}

// Enable mouse integration feature
func (engine *Akevitt) UseMouse() *Akevitt {
	engine.mouse = true

	return engine
}

// Register command with an alias and function
func (engine *Akevitt) UseRegisterCommand(command string, function CommandFunc) *Akevitt {
	command = strings.TrimSpace(command)
	engine.commands[command] = function
	return engine
}

// Engine default constructor
func NewEngine() *Akevitt {
	engine := &Akevitt{}
	engine.rooms = make(map[uint64]Room)
	engine.sessions = make(Sessions)
	engine.commands = make(map[string]CommandFunc)
	engine.bind = ":2222"
	engine.dbPath = "data/database.db"
	engine.mouse = false
	engine.heartbeats = make(map[int]chan struct{})
	return engine
}

// Sets the spawn room.
// Note: During startup, the engine traverses from spawn room to exits associated with that room recursively.
// Make sure you connect rooms with BindRoom function
func (engine *Akevitt) UseSpawnRoom(r Room) *Akevitt {
	engine.defaultRoom = r

	return engine
}

func (engine *Akevitt) UseNewHeartbeat(interval int) *Akevitt {
	engine.heartbeats[interval] = make(chan struct{})
	return engine
}
