package akevitt

import (
	"io"

	"github.com/rivo/tview"
)

type ActiveSession struct {
	Account *Account
	UI      *tview.Application
	prevUI  *tview.Primitive
}

func (session *ActiveSession) SetRoot(primitive tview.Primitive) {
	session.prevUI = &primitive
	session.UI.SetRoot(primitive, false)
}

func (session *ActiveSession) GetPreviousUI() *tview.Primitive {
	return session.prevUI
}

func purgeDeadSessions(sessions *Sessions) {
	for k := range *sessions {
		_, err := io.WriteString(k, "\000")
		if err != nil {
			delete(*sessions, k)
		}
	}
}
