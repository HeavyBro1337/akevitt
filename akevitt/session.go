package akevitt

import (
	"io"

	"github.com/rivo/tview"
)

type ActiveSession struct {
	Account *Account
	UI      *tview.Application
}

func purgeDeadSessions(sessions *Sessions) {
	for k := range *sessions {
		_, err := io.WriteString(k, "\000")
		if err != nil {
			delete(*sessions, k)
		}
	}
}
