package akevitt

import (
	"fmt"

	"github.com/gliderlabs/ssh"
)

// Broadcasts message
func broadcastMessage(sessions map[ssh.Session]*ActiveSession, message string, sender *ActiveSession,
	onMessage func(message string, sender *ActiveSession, currentSession *ActiveSession)) error {
	fmt.Printf("sender: %v\n", sender.Account)
	fmt.Printf("len(sessions): %v\n", len(sessions))
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
