package plugins

import (
	"fmt"

	"github.com/IvanKorchmit/akevitt"
	"github.com/rivo/tview"
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

func (plugin *MessagePlugin) UpdateChannel(old, new string, session *akevitt.ActiveSession) {
	channels := session.Data[MessagePluginData].([]string)
	channels = akevitt.RemoveItem(channels, old)
	channels = append(channels, new)
	session.Data[MessagePluginData] = channels
}

func (plugin *MessagePlugin) Build(engine *akevitt.Akevitt) error {
	if plugin.includeCmd {
		engine.AddCommand("ooc", plugin.oocCmd)
	}
	engine.AddInit(func(session *akevitt.ActiveSession) {
		textView := tview.NewTextView()
		session.Data[logElem] = textView
		session.Data[MessagePluginData] = []string{"ooc"}
	})
	return nil
}

func NewMessagePlugin(includeCmd bool, fn MessageFunc, format string) *MessagePlugin {
	if format == "" {
		format = "%s (%s) says %s"
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

func (plugin *MessagePlugin) GetChatLog(session *akevitt.ActiveSession) *tview.TextView {
	tv := session.Data[logElem].(*tview.TextView)

	return tv
}
