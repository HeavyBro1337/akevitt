package akevitt

import (
	"github.com/gliderlabs/ssh"
)

// Broadcasts message
func broadcastMessage(sessions map[ssh.Session]*ActiveSession, message string, sender *ActiveSession,
	onMessage func(message string, sender *ActiveSession, currentSession *ActiveSession)) error {
	for _, element := range sessions {
		// The user is not authenticated
		if element.Account == nil {
			continue
		}
		onMessage(message, sender, element)
		// element.ui.Draw()
		if element != sender {
			element.UI.ForceDraw()
		}
	}
	return nil
}
