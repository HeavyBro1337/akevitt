package akevitt

import (
	"io"

	"github.com/rivo/tview"
)

type ActiveSession interface {
	GetAccount() *Account
	SetAccount(acc *Account)
	GetApplication() *tview.Application
	SetApplication(app *tview.Application)
}

func purgeDeadSessions(sessions *Sessions, engine *Akevitt, callback DeadSessionFunc) {
	deadSessions := make([]ActiveSession, 0)
	liveSessions := make([]ActiveSession, 0)
	for k, v := range *sessions {
		_, err := io.WriteString(k, "\000")
		if err != nil {
			delete(*sessions, k)
			deadSessions = append(deadSessions, v)
			continue
		}
		liveSessions = append(liveSessions, v)
	}

	for _, v := range deadSessions {
		callback(v, liveSessions, engine)
	}
}
