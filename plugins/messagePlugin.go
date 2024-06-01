package plugins

import (
	"fmt"

	"github.com/IvanKorchmit/akevitt"
	"github.com/rivo/tview"
)

type mapLogs = map[string]*tview.TextView

type MessageFunc = func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, channel, message, username string) error

type MessagePlugin struct {
	onMessageFn MessageFunc
	includeCmd  bool
	format      string
	sessions    map[*akevitt.ActiveSession]mapLogs
}

// Send the message to other current sessions
func (plugin *MessagePlugin) Message(engine *akevitt.Akevitt, channel, message, username string, session *akevitt.ActiveSession) error {
	akevitt.PurgeDeadSessions(engine, engine.GetOnDeadSession()...)

	for _, v := range engine.GetSessions() {
		tvChannel, ok := plugin.sessions[v][channel]

		if !ok {
			continue
		}

		tvAll := plugin.sessions[v]["all"]

		if plugin.onMessageFn != nil {
			err := plugin.onMessageFn(engine, v, channel, message, username)

			if err != nil {
				return err
			}
		}

		st := fmt.Sprintf(plugin.format, username, channel, message)

		akevitt.AppendText(st, tvChannel)
		akevitt.AppendText(st, tvAll)

		if session != v {
			v.Application.Draw()
		}
	}

	return nil
}

func (plugin *MessagePlugin) UpdateChannel(old, new string, session *akevitt.ActiveSession) {
	plugin.sessions[session][old] = nil
	delete(plugin.sessions[session], old)
	plugin.sessions[session][new] = tview.NewTextView()
}

func (plugin *MessagePlugin) GetChannels(session *akevitt.ActiveSession) []string {
	return akevitt.GetMapKeys(plugin.sessions[session])
}

func (plugin *MessagePlugin) AddChannel(channel string, session *akevitt.ActiveSession) {
	_, ok := plugin.sessions[session][channel]

	if ok {
		return
	}

	plugin.sessions[session][channel] = tview.NewTextView()
}

func (plugin *MessagePlugin) RemoveChannel(channel string, session *akevitt.ActiveSession) {
	delete(plugin.sessions[session], channel)
}

func (plugin *MessagePlugin) Build(engine *akevitt.Akevitt) error {
	if plugin.includeCmd {
		engine.AddCommand("ooc", plugin.oocCmd)
	}
	engine.AddInit(func(engine *akevitt.Akevitt, session *akevitt.ActiveSession) {
		plugin.sessions[session] = make(map[string]*tview.TextView)
		plugin.sessions[session]["all"] = tview.NewTextView()
		plugin.sessions[session]["ooc"] = tview.NewTextView()
	})

	engine.AddSessionDead(func(deadSession *akevitt.ActiveSession, liveSessions []*akevitt.ActiveSession, engine *akevitt.Akevitt) {
		delete(plugin.sessions, deadSession)
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
		sessions:    make(map[*akevitt.ActiveSession]map[string]*tview.TextView),
	}
}

// Out-of-character chat command
func (plugin *MessagePlugin) oocCmd(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {

	return plugin.Message(engine, "ooc", command, session.Account.Username, session)
}

func (plugin *MessagePlugin) GetChatLog(session *akevitt.ActiveSession) *tview.TextView {
	tv := plugin.sessions[session]["all"]

	return tv
}

func (plugin *MessagePlugin) GetChatLogChannel(channel string, session *akevitt.ActiveSession) *tview.TextView {
	tv := plugin.sessions[session][channel]

	return tv
}
