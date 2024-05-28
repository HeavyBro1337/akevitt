package plugins

import (
	"fmt"
	"strings"

	"github.com/IvanKorchmit/akevitt"
)

type Doc struct {
	Long  string
	Brief string
}

type HelpPlugin struct {
	docs map[string]Doc
}

type helpPluginBuilder struct {
	plugin *HelpPlugin
}

func NewHelpPlugin() *helpPluginBuilder {
	return &helpPluginBuilder{
		plugin: &HelpPlugin{
			docs: make(map[string]Doc),
		},
	}
}

func (plugin *HelpPlugin) Build(engine *akevitt.Akevitt) error {
	lenDocs := len(plugin.docs)
	lenCommands := len(engine.GetCommands())

	if lenDocs != lenCommands {
		return fmt.Errorf("help plugin: amount of documented commands are %d, but the game has %d", lenDocs, lenCommands)
	}

	for _, cmd := range engine.GetCommands() {
		_, ok := plugin.docs[cmd]

		if !ok {
			return fmt.Errorf("help plugin: the command %s is undocumented", cmd)
		}
	}

	return nil
}

func (plugin *HelpPlugin) FindHelp(command string) (string, error) {
	doc, ok := plugin.docs[command]

	if !ok {
		return "", fmt.Errorf("help for the command \"%s\" not found", command)
	}

	return doc.Long, nil
}

func (plugin *HelpPlugin) ListHelp() string {
	format := "%s --- %s\n"

	builder := strings.Builder{}
	for cmd, doc := range plugin.docs {
		builder.WriteString(fmt.Sprintf(format, cmd, doc.Brief))
	}
	return builder.String()
}

func (plugin *HelpPlugin) ListHelpCustom(callback func(string, Doc) string) string {
	builder := strings.Builder{}

	for cmd, v := range plugin.docs {
		builder.WriteString(callback(cmd, v))
	}

	return builder.String()
}
