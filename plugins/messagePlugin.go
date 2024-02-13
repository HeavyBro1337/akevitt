package plugins

import (
	"errors"

	"github.com/IvanKorchmit/akevitt"
)

const MessagePluginData string = "MessagePlugin"

type MessageFunc = func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, channel, message, username string) error

type MessagePlugin struct {
	onMessageFn MessageFunc
	includeCmd  bool
}

// Send the message to other current sessions
func (plugin *MessagePlugin) Message(engine *akevitt.Akevitt, channel, message, username string, session *akevitt.ActiveSession) error {
	if plugin.onMessageFn == nil {
		return errors.New("message callback is nil")
	}
	akevitt.PurgeDeadSessions(engine, engine.GetOnDeadSession())

	for _, v := range engine.GetSessions() {
		channels := v.Data[MessagePluginData].([]string)

		if !akevitt.Find(channels, channel) {
			continue
		}

		err := plugin.onMessageFn(engine, v, channel, message, username)

		if session != v {
			v.Application.Draw()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (plugin *MessagePlugin) Build(engine *akevitt.Akevitt) error {
	if plugin.includeCmd {
		engine.AddCommand("ooc", plugin.oocCmd)
	}
	engine.AddInit(func(session *akevitt.ActiveSession) {
		session.Data[MessagePluginData] = []string{"ooc"}
	})
	return nil
}

func NewMessagePlugin(includeCmd bool, fn MessageFunc) *MessagePlugin {
	return &MessagePlugin{
		includeCmd:  includeCmd,
		onMessageFn: fn,
	}
}

// Out-of-character chat command
func (plugin *MessagePlugin) oocCmd(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {

	return plugin.Message(engine, "ooc", command, session.Account.Username, session)
}
