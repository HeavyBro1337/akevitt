package network

import (
	"akevitt/core/database/credentials"

	"github.com/gliderlabs/ssh"
)

// Broadcasts message
func BroadcastMessage(sessions *map[ssh.Session]ActiveSession, message string, session ssh.Session,
	onMessage func(message string, sender credentials.Account, currentSession ActiveSession)) error {
	for key, element := range *sessions {
		// The user is not authenticated
		if element.Account == nil {
			continue
		}
		onMessage(message, *(*sessions)[session].Account, element)
		// element.ui.Draw()
		if key != session {
			element.UI.ForceDraw()
		}
	}
	return nil
}
