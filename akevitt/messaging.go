package akevitt

import (
	"github.com/gliderlabs/ssh"
)

// Broadcasts message
func broadcastMessage(sessions map[ssh.Session]*ActiveSession, message string, sender *ActiveSession,
	onMessage func(message string, sender *ActiveSession, currentSession *ActiveSession)) error {
	for _, session := range sessions {
		// The user is not authenticated
		if session.Account == nil {
			continue
		}
		onMessage(message, sender, session)
		// element.ui.Draw()
		if session != sender {
			session.UI.ForceDraw()
		}
	}
	return nil
}
