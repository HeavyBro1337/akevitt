package engine

import (
	"strings"
	"time"
)

func (builder *akevittBuilder) UseBind(bindAddress string) *akevittBuilder {
	builder.engine.bind = bindAddress

	return builder
}

func (builder *akevittBuilder) UseRootUI(uiFunc UIFunc) *akevittBuilder {
	builder.engine.root = uiFunc

	return builder
}

func (builder *akevittBuilder) UseRegisterCommand(command string, function CommandFunc) *akevittBuilder {
	builder.engine.AddCommand(command, function)
	return builder
}

func (engine *Akevitt) AddInit(fn func(*Akevitt, *ActiveSession)) {
	engine.initFunc = append(engine.initFunc, fn)
}

func (engine *Akevitt) AddCommand(command string, function CommandFunc) {
	command = strings.TrimSpace(command)
	engine.commands[command] = function
}

func NewEngine() *akevittBuilder {
	engine := &Akevitt{}
	engine.rooms = make(map[uint64]*Room)
	engine.roomsByName = make(map[string]*Room)
	engine.roomsByGUID = make(map[string]*Room)
	engine.npcs = make(map[string]*NPC)
	engine.items = make(map[string]*Item)
	engine.sessions = make(Sessions)
	engine.commands = make(map[string]CommandFunc)
	engine.bind = ":1999"
	engine.rsaKey = "id_rsa"
	engine.mouse = false
	engine.heartbeats = make(map[int]*Pair[time.Ticker, []func() error])
	engine.plugins = make([]Plugin, 0)
	engine.luaVM = NewLuaVM()

	builder := &akevittBuilder{engine}

	return builder
}

func (builder *akevittBuilder) UseSpawnRoom(r *Room) *akevittBuilder {
	builder.engine.defaultRoom = r

	return builder
}

func (builder *akevittBuilder) UseOnJoin(fn func(*Akevitt, *ActiveSession)) *akevittBuilder {
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

func (engine *Akevitt) AddRoom(room *Room) error {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	if room.GUID == "" {
		room.GUID = generateGUID()
	}

	engine.rooms[room.GetKey()] = room
	engine.roomsByName[room.Name] = room
	engine.roomsByGUID[room.GUID] = room

	return nil
}

func (engine *Akevitt) GetRoomByName(name string) *Room {
	engine.mu.RLock()
	defer engine.mu.RUnlock()
	return engine.roomsByName[name]
}

func (engine *Akevitt) GetRoomByGUID(guid string) *Room {
	engine.mu.RLock()
	defer engine.mu.RUnlock()
	return engine.roomsByGUID[guid]
}

func (engine *Akevitt) AddNPC(npc *NPC) {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	if npc.GUID == "" {
		npc.GUID = generateGUID()
	}

	engine.npcs[npc.GUID] = npc
}

func (engine *Akevitt) GetNPCByGUID(guid string) *NPC {
	engine.mu.RLock()
	defer engine.mu.RUnlock()
	return engine.npcs[guid]
}

func (engine *Akevitt) GetNPCByName(name string) *NPC {
	engine.mu.RLock()
	defer engine.mu.RUnlock()
	for _, npc := range engine.npcs {
		if npc.Name == name {
			return npc
		}
	}
	return nil
}

func (engine *Akevitt) AddItem(item *Item) {
	engine.mu.Lock()
	defer engine.mu.Unlock()

	if item.GUID == "" {
		item.GUID = generateGUID()
	}

	engine.items[item.GUID] = item
}

func (engine *Akevitt) GetItemByGUID(guid string) *Item {
	engine.mu.RLock()
	defer engine.mu.RUnlock()
	return engine.items[guid]
}

func (engine *Akevitt) GetLuaVM() *LuaVM {
	return engine.luaVM
}

func (engine *Akevitt) LoadLuaScript(path string) error {
	engine.luaVM.SetEngine(engine)
	return engine.luaVM.LoadScript(path)
}

func (engine *Akevitt) LoadLuaScriptsDir(dir string) error {
	engine.luaVM.SetEngine(engine)
	return engine.luaVM.LoadScriptsDir(dir)
}

func (engine *Akevitt) ExecuteLuaCommand(command string, session *ActiveSession) error {
	engine.luaVM.SetEngine(engine)
	return engine.luaVM.CallCommand(command, session)
}

func (engine *Akevitt) ReloadLua() error {
	if err := engine.luaVM.Reload(); err != nil {
		return err
	}
	engine.luaVM.SetEngine(engine)
	return nil
}