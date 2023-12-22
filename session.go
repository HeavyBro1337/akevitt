package akevitt

import (
	"io"

	"github.com/rivo/tview"
)

type ActiveSession struct {
	Account     *Account
	Application *tview.Application
	Data        map[string]any
}

func purgeDeadSessions(sessions *Sessions, engine *Akevitt, callback DeadSessionFunc) {
	deadSessions := make([]*ActiveSession, 0)
	liveSessions := make([]*ActiveSession, 0)
	for k, v := range *sessions {
		_, err := io.WriteString(k, "\000")
		if err != nil {
			delete(*sessions, k)
			deadSessions = append(deadSessions, v)
			continue
		}
		liveSessions = append(liveSessions, v)
	}
	if callback != nil {
		for _, v := range deadSessions {
			callback(v, liveSessions, engine)
		}
	}
}
