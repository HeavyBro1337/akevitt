package basic

import (
	"errors"
	"fmt"

	"github.com/IvanKorchmit/akevitt"
)

func OnMessageHook(engine *akevitt.Akevitt, session *Session, channel, message, username string) error {
	if session == nil {
		return errors.New("session is nil. Probably the dead one")
	}

	st := fmt.Sprintf("%s (%s): %s", username, channel, message)

	if session.subscribedChannels != nil {
		if akevitt.Find[string](session.subscribedChannels, channel) {
			AppendText(session, st, session.Chat)
		} else if session.Character.currentRoom.GetName() == channel {
			AppendText(session, st, session.Chat)
		}
	} else {
		fmt.Printf("warn: the channels is empty at %s", session.account.Username)
	}

	return nil
}

// In this hook it is good to do some clean up (I.e. removing associated character from a room)
func OnSessionDeathHook(deadSession *Session, liveSessions []*Session, engine *akevitt.Akevitt) {
	if deadSession.account == nil {
		return
	}

	deadSession.Character.currentRoom.RemoveObject(deadSession.Character)
	for _, v := range liveSessions {
		if v.Chat == nil {
			continue
		}

		AppendText(v, fmt.Sprintf("%s left the game", deadSession.account.Username), v.Chat)
	}
}

func OnDialogue(engine *akevitt.Akevitt, session Session, dialogue *akevitt.Dialogue) error {
	return DialogueBox(dialogue, engine, &session)
}
