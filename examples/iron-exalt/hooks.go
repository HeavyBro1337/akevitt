package main

import (
	"errors"
	"fmt"

	"github.com/IvanKorchmit/akevitt"
)

func onMessage(engine *akevitt.Akevitt, session akevitt.ActiveSession, channel, message, username string) error {
	if session == nil {
		return errors.New("session is nil. Probably the dead one")
	}

	sess, ok := session.(*IronExaltSession)

	st := fmt.Sprintf("%s (%s): %s", username, channel, message)

	if ok && sess.subscribedChannels != nil {
		if akevitt.Find[string](sess.subscribedChannels, channel) {
			AppendText(sess, st, sess.chat)
		} else if sess.character.currentRoom.GetName() == channel {
			AppendText(sess, st, sess.chat)
		}
	} else if !ok {
		fmt.Printf("could not cast to session")
	}

	return nil
}

// In this hook it is good to do some clean up (I.e. removing associated character from a room)
func onSessionEnd(deadSession akevitt.ActiveSession, liveSessions []akevitt.ActiveSession, engine *akevitt.Akevitt) {
	sess, ok := deadSession.(*IronExaltSession)
	if !ok {
		fmt.Println("could not cast to session")
		return
	}
	if sess.account == nil {
		return
	}

	sess.character.currentRoom.RemoveObject(sess.character)
	for _, v := range liveSessions {
		lsess, ok := v.(*IronExaltSession)

		if !ok || lsess.chat == nil {
			continue
		}

		AppendText(lsess, fmt.Sprintf("%s left the game", sess.account.Username), lsess.chat)
	}
}

func onDialogue(engine *akevitt.Akevitt, session akevitt.ActiveSession, dialogue *akevitt.Dialogue) error {
	sess, ok := session.(*IronExaltSession)

	if !ok {
		return errors.New("could not cast to session")
	}

	err := dialogueBox(dialogue, engine, sess)
	return err
}
