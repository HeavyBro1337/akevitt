package akevitt

import "strings"

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
	return engine
}

func (engine *Akevitt) UseSpawnRoom(r Room) *Akevitt {
	engine.defaultRoom = r

	return engine
}
