package plugins

import (
	"errors"
	"fmt"

	"github.com/IvanKorchmit/akevitt"
	"github.com/uaraven/logview"
)

const MessagePluginData string = "MessagePlugin"
const logElem string = "MessagePluginLog"

type MessageFunc = func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, channel, message, username string) error

type MessagePlugin struct {
	onMessageFn MessageFunc
	includeCmd  bool
	format      string
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
		if plugin.onMessageFn != nil {
			err := plugin.onMessageFn(engine, v, channel, message, username)

			if err != nil {
				return err
			}
		}

		st := fmt.Sprintf(plugin.format, username, channel, message)

		akevitt.AppendText(st, plugin.GetChatLog(v))

		if session != v {
			v.Application.Draw()
		}
	}

	return nil
}

func (plugin *MessagePlugin) Build(engine *akevitt.Akevitt) error {
	if plugin.includeCmd {
		engine.AddCommand("ooc", plugin.oocCmd)
	}
	engine.AddInit(func(session *akevitt.ActiveSession) {
		session.Data[logElem] = logview.NewLogView()
		session.Data[MessagePluginData] = []string{"ooc"}
	})
	return nil
}

func NewMessagePlugin(includeCmd bool, fn MessageFunc, format string) *MessagePlugin {
	if format == "" {
		format = "%[0]s (%[1]s) says %[2]s"
	}

	return &MessagePlugin{
		includeCmd:  includeCmd,
		onMessageFn: fn,
		format:      format,
	}
}

// Out-of-character chat command
func (plugin *MessagePlugin) oocCmd(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {

	return plugin.Message(engine, "ooc", command, session.Account.Username, session)
}

func (plugin *MessagePlugin) GetChatLog(session *akevitt.ActiveSession) *logview.LogView {
	lv := session.Data[logElem].(*logview.LogView)

	return lv
}
