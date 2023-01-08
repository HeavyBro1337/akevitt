package network

import (
	"akevitt/core/database/credentials"
	"io"

	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type ActiveSession struct {
	Account *credentials.Account
	UI      *tview.Application
	Chat    *tview.List
}

// Iterates through all currently dead sessions by trying to send null character.
// If it gets an error, then we found the dead session and we purge them from active ones.
func PurgeDeadSessions(sessions *map[ssh.Session]ActiveSession) {
	for k := range *sessions {

		_, err := io.WriteString(k, "\000")
		if err != nil {
			delete(*sessions, k)
		}
	}
}
