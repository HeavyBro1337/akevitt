package engine

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/yuin/gopher-lua"
)

type LuaVM struct {
	mu        sync.Mutex
	L         *lua.LState
	scripts   map[string]*lua.LFunction
	scriptsPath string
}

func NewLuaVM() *LuaVM {
	vm := &LuaVM{
		scripts: make(map[string]*lua.LFunction),
	}
	vm.L = lua.NewState()
	return vm
}

func (vm *LuaVM) SetEngine(engine *Akevitt) {
	vm.registerAPI(engine)
}

func (vm *LuaVM) SetScriptsPath(path string) {
	vm.scriptsPath = path
}

func (vm *LuaVM) LoadScriptsDir(dir string) error {
	entries, err := filepath.Glob(filepath.Join(dir, "*.lua"))
	if err != nil {
		return fmt.Errorf("failed to glob scripts: %w", err)
	}

	for _, path := range entries {
		if err := vm.LoadScript(path); err != nil {
			return err
		}
	}

	return nil
}

func (vm *LuaVM) registerAPI(engine *Akevitt) {
	vm.L.PreloadModule("akevitt", func(L *lua.LState) int {
		mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"version":           vm.apiVersion,
			"addRoom":           vm.makeRoomFunc(engine),
			"getRoom":           vm.getRoomFunc(engine),
			"bindRooms":         vm.bindRoomsFunc(engine),
			"addNPC":            vm.addNPCFunc(engine),
			"sendMessage":       vm.sendMessageFunc(engine),
			"sendMessageToRoom": vm.sendMessageToRoomFunc(engine),
			"getSessions":       vm.getSessionsFunc(engine),
			"getPlayerRoom":     vm.getPlayerRoomFunc(engine),
			"setPlayerRoom":     vm.setPlayerRoomFunc(engine),
			"setPlayerData":     vm.setPlayerDataFunc(engine),
			"getPlayerData":     vm.getPlayerDataFunc(engine),
			"createItem":        vm.createItemFunc(engine),
			"getRoomExits":      vm.getRoomExitsFunc(engine),
		})
		L.Push(mod)
		return 1
	})
}

func (vm *LuaVM) apiVersion(L *lua.LState) int {
	L.Push(lua.LString("0.1.0"))
	return 1
}

func (vm *LuaVM) makeRoomFunc(engine *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		opts := L.CheckTable(1)
		name := opts.RawGetString("name")
		description := opts.RawGetString("description")

		if name.Type() != lua.LTString {
			L.RaiseError("room name must be a string")
			return 0
		}

		desc := ""
		if description.Type() == lua.LTString {
			desc = description.String()
		}

		room := &Room{
			ObjectImpl: ObjectImpl{
				Name: name.String(),
			},
			Exits:   make([]*Exit, 0),
			Objects: make([]Object, 0),
		}

		if desc != "" {
			room.Description = desc
		}

		if err := engine.AddRoom(room); err != nil {
			L.RaiseError("failed to add room: %v", err)
			return 0
		}

		t := L.NewTable()
		t.RawSetString("guid", lua.LString(room.GUID))
		t.RawSetString("name", lua.LString(room.Name))
		L.Push(t)
		return 1
	}
}

func (vm *LuaVM) getRoomFunc(engine *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		name := L.CheckString(1)

		room := engine.GetRoomByName(name)
		if room == nil {
			L.Push(lua.LNil)
			return 1
		}

		t := vm.roomToTable(room)
		L.Push(t)
		return 1
	}
}

func (vm *LuaVM) bindRoomsFunc(engine *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		fromName := L.CheckString(1)
		toName := L.CheckString(2)

		opts := L.OptTable(3, nil)
		oneWay := false
		if opts != nil {
			if oneWayVal := opts.RawGetString("one_way"); oneWayVal.Type() == lua.LTBool {
				oneWay = bool(oneWayVal.(lua.LBool))
			}
		}

		fromRoom := engine.GetRoomByName(fromName)
		toRoom := engine.GetRoomByName(toName)

		if fromRoom == nil {
			L.RaiseError("room '%s' not found", fromName)
			return 0
		}
		if toRoom == nil {
			L.RaiseError("room '%s' not found", toName)
			return 0
		}

		exit := Exit{Room: toRoom}
		BindRooms(fromRoom, exit, toRoom)

		if !oneWay {
			exitBack := Exit{Room: fromRoom}
			BindRooms(toRoom, exitBack, fromRoom)
		}

		return 0
	}
}

func (vm *LuaVM) addNPCFunc(engine *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		opts := L.CheckTable(1)
		name := opts.RawGetString("name")
		description := opts.RawGetString("description")
		roomName := opts.RawGetString("room")

		if name.Type() != lua.LTString {
			L.RaiseError("npc name must be a string")
			return 0
		}

		npc := &NPC{
			ObjectImpl: ObjectImpl{
				Name:        name.String(),
				Description: description.String(),
			},
		}

		if roomName.Type() == lua.LTString {
			room := engine.GetRoomByName(roomName.String())
			if room != nil {
				npc.RoomID = room.GUID
				room.Objects = append(room.Objects, npc)
			}
		}

		engine.AddNPC(npc)

		t := L.NewTable()
		t.RawSetString("guid", lua.LString(npc.GUID))
		t.RawSetString("name", lua.LString(npc.Name))
		L.Push(t)
		return 1
	}
}

func (vm *LuaVM) sendMessageFunc(_ *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		session := L.CheckUserData(1)
		message := L.CheckString(2)

		activeSession, ok := session.Value.(*ActiveSession)
		if !ok {
			L.RaiseError("invalid session")
			return 0
		}

		activeSession.Send(message)
		return 0
	}
}

func (vm *LuaVM) sendMessageToRoomFunc(engine *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		roomName := L.CheckString(1)
		message := L.CheckString(2)
		excludeSession := L.OptUserData(3, nil)

		room := engine.GetRoomByName(roomName)
		if room == nil {
			L.RaiseError("room '%s' not found", roomName)
			return 0
		}

		for _, s := range engine.sessions {
			if excludeSession != nil {
				ex, _ := excludeSession.Value.(*ActiveSession)
				if s == ex {
					continue
				}
			}
			if s.RoomID == room.GUID {
				s.Send(message)
			}
		}

		return 0
	}
}

func (vm *LuaVM) getSessionsFunc(engine *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		t := L.NewTable()
		i := 1
		for _, s := range engine.sessions {
			st := vm.sessionToTable(s)
			t.RawSetInt(i, st)
			i++
		}
		L.Push(t)
		return 1
	}
}

func (vm *LuaVM) getPlayerRoomFunc(engine *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		session := L.CheckUserData(1)
		activeSession, ok := session.Value.(*ActiveSession)
		if !ok {
			L.RaiseError("invalid session")
			return 0
		}

		room := engine.GetRoomByGUID(activeSession.RoomID)
		if room == nil {
			L.Push(lua.LNil)
			return 1
		}

		t := vm.roomToTable(room)
		L.Push(t)
		return 1
	}
}

func (vm *LuaVM) setPlayerRoomFunc(engine *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		session := L.CheckUserData(1)
		roomName := L.CheckString(2)

		activeSession, ok := session.Value.(*ActiveSession)
		if !ok {
			L.RaiseError("invalid session")
			return 0
		}

		room := engine.GetRoomByName(roomName)
		if room == nil {
			L.RaiseError("room '%s' not found", roomName)
			return 0
		}

		activeSession.RoomID = room.GUID
		return 0
	}
}

func (vm *LuaVM) setPlayerDataFunc(_ *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		session := L.CheckUserData(1)
		key := L.CheckString(2)
		value := L.CheckAny(3)

		activeSession, ok := session.Value.(*ActiveSession)
		if !ok {
			L.RaiseError("invalid session")
			return 0
		}

		activeSession.Data[key] = vm.lValueToAny(value)
		return 0
	}
}

func (vm *LuaVM) getPlayerDataFunc(_ *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		session := L.CheckUserData(1)
		key := L.CheckString(2)

		activeSession, ok := session.Value.(*ActiveSession)
		if !ok {
			L.RaiseError("invalid session")
			return 0
		}

		val, ok := activeSession.Data[key]
		if !ok {
			L.Push(lua.LNil)
			return 1
		}

		L.Push(vm.anyToLValue(val))
		return 1
	}
}

func (vm *LuaVM) createItemFunc(engine *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		opts := L.CheckTable(1)
		name := opts.RawGetString("name")
		description := opts.RawGetString("description")

		if name.Type() != lua.LTString {
			L.RaiseError("item name must be a string")
			return 0
		}

		desc := ""
		if description.Type() == lua.LTString {
			desc = description.String()
		}

		item := &Item{
			ObjectImpl: ObjectImpl{
				Name:        name.String(),
				Description: desc,
			},
		}

		engine.AddItem(item)

		t := L.NewTable()
		t.RawSetString("guid", lua.LString(item.GUID))
		t.RawSetString("name", lua.LString(item.Name))
		L.Push(t)
		return 1
	}
}

func (vm *LuaVM) getRoomExitsFunc(engine *Akevitt) lua.LGFunction {
	return func(L *lua.LState) int {
		roomName := L.CheckString(1)

		room := engine.GetRoomByName(roomName)
		if room == nil {
			L.RaiseError("room '%s' not found", roomName)
			return 0
		}

		t := L.NewTable()
		for i, exit := range room.Exits {
			et := L.NewTable()
			et.RawSetString("target", lua.LString(exit.Room.Name))
			et.RawSetString("target_guid", lua.LString(exit.Room.GUID))
			t.RawSetInt(i+1, et)
		}

		L.Push(t)
		return 1
	}
}

func (vm *LuaVM) roomToTable(room *Room) *lua.LTable {
	t := vm.L.NewTable()
	t.RawSetString("guid", lua.LString(room.GUID))
	t.RawSetString("name", lua.LString(room.Name))
	t.RawSetString("description", lua.LString(room.Description))

	exits := vm.L.NewTable()
	for i, exit := range room.Exits {
		et := vm.L.NewTable()
		et.RawSetString("target", lua.LString(exit.Room.Name))
		et.RawSetString("target_guid", lua.LString(exit.Room.GUID))
		exits.RawSetInt(i+1, et)
	}
	t.RawSetString("exits", exits)

	objects := vm.L.NewTable()
	for i, obj := range room.Objects {
		ot := vm.L.NewTable()
		ot.RawSetString("name", lua.LString(obj.GetName()))
		objects.RawSetInt(i+1, ot)
	}
	t.RawSetString("objects", objects)

	return t
}

func (vm *LuaVM) sessionToTable(session *ActiveSession) *lua.LTable {
	t := vm.L.NewTable()

	if session.Account != nil {
		acc := vm.L.NewTable()
		acc.RawSetString("username", lua.LString(session.Account.Username))
		t.RawSetString("account", acc)
	}

	t.RawSetString("room_id", lua.LString(session.RoomID))
	t.RawSetString("data", vm.anyToTable(session.Data))

	t.RawSetString("__userdata", &lua.LUserData{Value: session})
	return t
}

func (vm *LuaVM) anyToTable(data map[string]any) *lua.LTable {
	t := vm.L.NewTable()
	for k, v := range data {
		t.RawSetString(k, vm.anyToLValue(v))
	}
	return t
}

func (vm *LuaVM) anyToLValue(v any) lua.LValue {
	switch val := v.(type) {
	case string:
		return lua.LString(val)
	case int:
		return lua.LNumber(val)
	case int64:
		return lua.LNumber(val)
	case float64:
		return lua.LNumber(val)
	case bool:
		return lua.LBool(val)
	case []string:
		t := vm.L.NewTable()
		for i, s := range val {
			t.RawSetInt(i+1, lua.LString(s))
		}
		return t
	case map[string]any:
		return vm.anyToTable(val)
	default:
		return lua.LNil
	}
}

func (vm *LuaVM) lValueToAny(v lua.LValue) any {
	switch val := v.(type) {
	case lua.LString:
		return string(val)
	case lua.LNumber:
		return float64(val)
	case lua.LBool:
		return bool(val)
	case *lua.LTable:
		return vm.tableToMap(val)
	default:
		return nil
	}
}

func (vm *LuaVM) tableToMap(t *lua.LTable) map[string]any {
	result := make(map[string]any)
	t.ForEach(func(k, v lua.LValue) {
		if key, ok := k.(lua.LString); ok {
			result[string(key)] = vm.lValueToAny(v)
		}
	})
	return result
}

func (vm *LuaVM) LoadScript(path string) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	scriptName := filepath.Base(path)
	if _, ok := vm.scripts[scriptName]; ok {
		return fmt.Errorf("script %s already loaded", scriptName)
	}

	fn := vm.L.DoFile(path)
	if fn != nil {
		return fmt.Errorf("failed to load script %s: %w", path, fn)
	}

	top := vm.L.Get(-1)
	if top.Type() != lua.LTFunction {
		return fmt.Errorf("script %s must return a function", path)
	}

	vm.scripts[scriptName] = top.(*lua.LFunction)
	return nil
}

func (vm *LuaVM) GetScript(name string) (*lua.LFunction, bool) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	fn, ok := vm.scripts[name]
	return fn, ok
}

func (vm *LuaVM) CallCommand(command string, session *ActiveSession) error {
	parts := SplitCommand(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	cmdName := parts[0]
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}

	sessionTable := vm.sessionToTable(session)
	cmdLua := lua.LString(cmdName)
	argsLua := lua.LString(args)

	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.L.SetGlobal("_cmd", &cmdLua)
	vm.L.SetGlobal("_args", &argsLua)
	vm.L.SetGlobal("_session", sessionTable)

	scriptName := cmdName + ".lua"
	fn, ok := vm.scripts[scriptName]
	if !ok {
		return fmt.Errorf("command %s not found", cmdName)
	}

	vm.L.Push(fn)
	vm.L.Push(sessionTable)
	vm.L.Push(&cmdLua)
	vm.L.Push(&argsLua)

	if err := vm.L.PCall(3, 0, nil); err != nil {
		return fmt.Errorf("error executing command %s: %w", cmdName, err)
	}

	return nil
}

func (vm *LuaVM) Reload() error {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.L.Close()
	vm.L = lua.NewState()
	vm.scripts = make(map[string]*lua.LFunction)

	return nil
}

func (vm *LuaVM) Close() {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	vm.L.Close()
}

func SplitCommand(command string) []string {
	if command == "" {
		return nil
	}

	inQuote := false
	quoteChar := byte(0)
	var result []string
	var current []byte

	for i := 0; i < len(command); i++ {
		c := command[i]

		if !inQuote && (c == '"' || c == '\'') {
			inQuote = true
			quoteChar = c
			continue
		}

		if inQuote && c == quoteChar {
			inQuote = false
			continue
		}

		if !inQuote && c == ' ' {
			if len(current) > 0 {
				result = append(result, string(current))
				current = nil
			}
			continue
		}

		current = append(current, c)
	}

	if len(current) > 0 {
		result = append(result, string(current))
	}

	return result
}

func GetScriptName(command string) string {
	parts := SplitCommand(command)
	if len(parts) == 0 {
		return ""
	}
	name := parts[0]
	name = strings.TrimSuffix(name, ".lua")
	return name + ".lua"
}