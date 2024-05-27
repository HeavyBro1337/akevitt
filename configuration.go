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

// Register command with an alias and function
func (builder *akevittBuilder) UseRegisterCommand(command string, function CommandFunc) *akevittBuilder {
	builder.engine.AddCommand(command, function)
	return builder
}

func (engine *Akevitt) AddInit(fn func(*ActiveSession)) {
	engine.initFunc = append(engine.initFunc, fn)
}

// Register command with an alias and function
func (engine *Akevitt) AddCommand(command string, function CommandFunc) {
	command = strings.TrimSpace(command)
	engine.commands[command] = function
}

// Engine default constructor
func NewEngine() *akevittBuilder {
	engine := &Akevitt{}
	engine.rooms = make(map[uint64]*Room)
	engine.sessions = make(Sessions)
	engine.commands = make(map[string]CommandFunc)
	engine.bind = ":2222"
	engine.rsaKey = "id_rsa"
	engine.dbPath = "data/database.db"
	engine.mouse = false
	engine.heartbeats = make(map[int]*Pair[time.Ticker, []func() error])
	engine.plugins = make([]Plugin, 0)

	builder := &akevittBuilder{engine}

	return builder
}

// Sets the spawn room.
// Note: During startup, the engine traverses from spawn room to exits associated with that room recursively.
// Make sure you connect rooms with BindRoom function
func (builder *akevittBuilder) UseSpawnRoom(r *Room) *akevittBuilder {
	builder.engine.defaultRoom = r

	return builder
}

func (builder *akevittBuilder) UseOnJoin(fn func(*ActiveSession)) *akevittBuilder {
	builder.engine.AddInit(fn)

	return builder
}

func (builder *akevittBuilder) UseKeyPath(path string) *akevittBuilder {
	builder.engine.rsaKey = path

	return builder
}

func (builder *akevittBuilder) AddPlugin(plugin ...Plugin) *akevittBuilder {
	builder.engine.addPlugin(plugin...)

	return builder
}

func (engine *Akevitt) addPlugin(plugins ...Plugin) {
	engine.plugins = append(engine.plugins, plugins...)
}
