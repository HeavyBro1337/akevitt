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

func purgeDeadSessions(sessions *Sessions) {
	for k := range *sessions {
		_, err := io.WriteString(k, "\000")
		if err != nil {
			println("Found dead")
			delete(*sessions, k)
		}
	}
}
