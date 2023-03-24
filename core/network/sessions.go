package network

import (
	"akevitt/core/objects/credentials"
	"io"

	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type ActiveSession struct {
	Account     *credentials.Account
	UI          *tview.Application
	Chat        *tview.List
	UIPrimitive *tview.Primitive
}

func (session ActiveSession) SetRoot(p tview.Primitive) {
	session.UIPrimitive = &p
	session.UI.SetRoot(p, true)
}

// Iterates through all current sessions by trying to send null character.
// If it receives an error, it indicates of session being dead.
func PurgeDeadSessions(sessions *map[ssh.Session]ActiveSession) {
	for k := range *sessions {

		_, err := io.WriteString(k, "\000")
		if err != nil {
			delete(*sessions, k)
		}
	}
}
