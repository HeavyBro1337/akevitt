package plugins

import (
	"fmt"

	"github.com/IvanKorchmit/akevitt/engine"
	"github.com/rivo/tview"
)

type mapLogs = map[string]*tview.TextView

type MessageFunc = func(engine *engine.Akevitt, session *engine.ActiveSession, channel, message, username string) error

type MessagePlugin struct {
	onMessageFn MessageFunc
	includeCmd  bool
	format      string
	sessions    map[*engine.ActiveSession]mapLogs
}

// Send the message to other current sessions
func (plugin *MessagePlugin) Message(eng *engine.Akevitt, channel, message, username string, session *engine.ActiveSession) error {
	engine.PurgeDeadSessions(eng, eng.GetOnDeadSession()...)

	for _, v := range eng.GetSessions() {
		tvChannel, ok := plugin.sessions[v][channel]

		if !ok {
			continue
		}

		tvAll := plugin.sessions[v]["all"]

		if plugin.onMessageFn != nil {
			err := plugin.onMessageFn(eng, v, channel, message, username)

			if err != nil {
				return err
			}
		}

		st := fmt.Sprintf(plugin.format, username, channel, message)

		engine.AppendText(st, tvChannel)
		engine.AppendText(st, tvAll)

		if session != v {
			v.Application.Draw()
		}
	}

	return nil
}

func (plugin *MessagePlugin) UpdateChannel(old, new string, session *engine.ActiveSession) {
	plugin.sessions[session][old] = nil
	delete(plugin.sessions[session], old)
	plugin.sessions[session][new] = tview.NewTextView()
}

func (plugin *MessagePlugin) GetChannels(session *engine.ActiveSession) []string {
	return engine.GetMapKeys(plugin.sessions[session])
}

func (plugin *MessagePlugin) AddChannel(channel string, session *engine.ActiveSession) {
	_, ok := plugin.sessions[session][channel]

	if ok {
		return
	}

	plugin.sessions[session][channel] = tview.NewTextView()
}

func (plugin *MessagePlugin) RemoveChannel(channel string, session *engine.ActiveSession) {
	delete(plugin.sessions[session], channel)
}

func (plugin *MessagePlugin) Build(eng *engine.Akevitt) error {
	if plugin.includeCmd {
		eng.AddCommand("ooc", plugin.oocCmd)
	}
	eng.AddInit(func(eng *engine.Akevitt, session *engine.ActiveSession) {
		plugin.sessions[session] = make(map[string]*tview.TextView)
		plugin.sessions[session]["all"] = tview.NewTextView()
		plugin.sessions[session]["ooc"] = tview.NewTextView()
	})

	eng.AddSessionDead(func(deadSession *engine.ActiveSession, liveSessions []*engine.ActiveSession, eng *engine.Akevitt) {
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
		sessions:    make(map[*engine.ActiveSession]map[string]*tview.TextView),
	}
}

// Out-of-character chat command
func (plugin *MessagePlugin) oocCmd(eng *engine.Akevitt, session *engine.ActiveSession, command string) error {

	return plugin.Message(eng, "ooc", command, session.Account.Username, session)
}

func (plugin *MessagePlugin) GetChatLog(session *engine.ActiveSession) *tview.TextView {
	tv := plugin.sessions[session]["all"]

	return tv
}

func (plugin *MessagePlugin) GetChatLogChannel(channel string, session *engine.ActiveSession) *tview.TextView {
	tv := plugin.sessions[session][channel]

	return tv
}