package plugins

import (
	"github.com/IvanKorchmit/akevitt/engine"
)

type LuaCommandPlugin struct {
	engine     *engine.Akevitt
	scriptsDir string
}

func NewLuaCommandPlugin(scriptsDir string) *LuaCommandPlugin {
	return &LuaCommandPlugin{
		scriptsDir: scriptsDir,
	}
}

func (plugin *LuaCommandPlugin) Build(eng *engine.Akevitt) error {
	plugin.engine = eng

	if err := eng.LoadLuaScriptsDir(plugin.scriptsDir); err != nil {
		return err
	}

	eng.AddCommand("*", plugin.handleLuaCommand)

	return nil
}

func (plugin *LuaCommandPlugin) handleLuaCommand(eng *engine.Akevitt, session *engine.ActiveSession, arguments string) error {
	return eng.ExecuteLuaCommand(arguments, session)
}