package akevitt

import (
	"strings"
	"time"
)

// Specify an address to listen.
// Example: :22, 127.0.0.1:2222, etc.
func (builder *akevittBuilder) UseBind(bindAddress string) *akevittBuilder {
	builder.engine.bind = bindAddress

	return builder
}

// Accepts function which returns the UI root screen.
func (builder *akevittBuilder) UseRootUI(uiFunc UIFunc) *akevittBuilder {
	builder.engine.root = uiFunc

	return builder
}

// Specify path to save database
func (builder *akevittBuilder) UseDBPath(path string) *akevittBuilder {
	builder.engine.dbPath = path

	return builder
}

// Enable mouse integration feature
func (builder *akevittBuilder) UseMouse() *akevittBuilder {
	builder.engine.mouse = true

	return builder
}

// Register command with an alias and function
func (builder *akevittBuilder) UseRegisterCommand(command string, function CommandFunc) *akevittBuilder {
	command = strings.TrimSpace(command)
	builder.engine.commands[command] = function
	return builder
}

// Engine default constructor
func NewEngine() *akevittBuilder {
	engine := &Akevitt{}
	engine.rooms = make(map[uint64]Room)
	engine.sessions = make(Sessions)
	engine.commands = make(map[string]CommandFunc)
	engine.bind = ":2222"
	engine.dbPath = "data/database.db"
	engine.mouse = false
	engine.heartbeats = make(map[int]*pair[time.Ticker, []func() error])

	builder := &akevittBuilder{engine}

	return builder
}

// Sets the spawn room.
// Note: During startup, the engine traverses from spawn room to exits associated with that room recursively.
// Make sure you connect rooms with BindRoom function
func (builder *akevittBuilder) UseSpawnRoom(r Room) *akevittBuilder {
	builder.engine.defaultRoom = r

	return builder
}

func (builder *akevittBuilder) UseNewHeartbeat(interval int) *akevittBuilder {
	dur := time.Duration(interval) * time.Second

	builder.engine.heartbeats[interval] = &pair[time.Ticker, []func() error]{f: *time.NewTicker(dur), s: make([]func() error, 0)}
	return builder
}

func (builder *akevittBuilder) UseOnJoin(f func(*ActiveSession)) *akevittBuilder {
	builder.engine.initFunc = f

	return builder
}
