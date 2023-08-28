package basic

import (
	"errors"
	"fmt"

	"github.com/IvanKorchmit/akevitt"
)

func OnMessageHook(engine *akevitt.Akevitt, session akevitt.ActiveSession, channel, message, username string) error {
	sess := CastSession[*Session](session)

	if session == nil {
		return errors.New("session is nil. Probably the dead one")
	}

	st := fmt.Sprintf("%s (%s): %s", username, channel, message)

	if sess.subscribedChannels != nil {
		if akevitt.Find[string](sess.subscribedChannels, channel) {
			AppendText(st, sess.Chat)
		} else if sess.Character.currentRoom.GetName() == channel {
			AppendText(st, sess.Chat)
		}
	} else {
		fmt.Printf("warn: the channels is empty at %s", sess.account.Username)
	}

	return nil
}

// In this hook it is good to do some clean up (I.e. removing associated character from a room)
func OnSessionDeathHook(deadSession akevitt.ActiveSession, liveSessions []akevitt.ActiveSession, engine *akevitt.Akevitt) {
	deadSess := CastSession[*Session](deadSession)

	if deadSess == nil {
		return
	}

	deadSess.Character.currentRoom.RemoveObject(deadSess.Character)
	for _, v := range akevitt.MapSlice(liveSessions, func(ls akevitt.ActiveSession) *Session {
		return CastSession[*Session](ls)
	}) {
		if v.Chat == nil {
			continue
		}

		AppendText(fmt.Sprintf("%s left the game", deadSess.account.Username), v.Chat)
	}
}

func OnDialogue(engine *akevitt.Akevitt, session akevitt.ActiveSession, dialogue *akevitt.Dialogue) error {
	sess := CastSession[*Session](session)

	return DialogueBox(dialogue, engine, sess)
}
