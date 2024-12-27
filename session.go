package akevitt

import (
	"io"

	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type ActiveSession struct {
	s           ssh.Session
	Application *tview.Application
	Data        map[string]any
}

func (a *ActiveSession) Kill(engine *Akevitt) error {
	return a.s.Exit(0)
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
