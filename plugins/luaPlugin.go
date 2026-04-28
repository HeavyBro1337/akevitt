package plugins

import (
	"github.com/IvanKorchmit/akevitt"
)

type LuaCommandPlugin struct {
	engine     *akevitt.Akevitt
	scriptsDir string
}

func NewLuaCommandPlugin(scriptsDir string) *LuaCommandPlugin {
	return &LuaCommandPlugin{
		scriptsDir: scriptsDir,
	}
}

func (plugin *LuaCommandPlugin) Build(engine *akevitt.Akevitt) error {
	plugin.engine = engine

	if err := engine.LoadLuaScriptsDir(plugin.scriptsDir); err != nil {
		return err
	}

	engine.AddCommand("*", plugin.handleLuaCommand)

	return nil
}

func (plugin *LuaCommandPlugin) handleLuaCommand(engine *akevitt.Akevitt, session *akevitt.ActiveSession, arguments string) error {
	return engine.ExecuteLuaCommand(arguments, session)
}
