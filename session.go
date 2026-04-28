package akevitt

import (
	"fmt"
	"io"
	"strings"

	"github.com/rivo/tview"
)

type ActiveSession struct {
	Account     *Account
	Application *tview.Application
	Data        map[string]any
	RoomID      string
}

func (session *ActiveSession) Send(message string) {
	fmt.Println(message)
}

func (session *ActiveSession) Sendf(format string, args ...any) {
	session.Send(fmt.Sprintf(format, args...))
}

func (session *ActiveSession) SendLines(lines ...string) {
	session.Send(strings.Join(lines, "\n"))
}

func (session *ActiveSession) loadData() {
	for k, v := range session.Account.PersistentData {
		session.Data[k] = v
	}
}

func PurgeDeadSessions(engine *Akevitt, callback ...DeadSessionFunc) {
	deadSessions := make([]*ActiveSession, 0)
	liveSessions := make([]*ActiveSession, 0)
	for k, v := range engine.sessions {
		_, err := io.WriteString(k, "\000")
		if err != nil {
			delete(engine.sessions, k)
			deadSessions = append(deadSessions, v)
			continue
		}
		liveSessions = append(liveSessions, v)
	}
	if callback != nil {
		for _, v := range deadSessions {
			for _, fn := range callback {
				fn(v, liveSessions, engine)
			}
		}
	}
}
