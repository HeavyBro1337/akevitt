package main

import (
	"akevitt/akevitt"
	"errors"
)

func say(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*ActiveSession)

	if !ok {
		return errors.New("invalid session type")
	}

	return engine.Message(sess.character.currentRoom.GetName(), arguments, sess.character.Name, session)
}

func ooc(engine *akevitt.Akevitt, session akevitt.ActiveSession, command string) error {
	return engine.Message("ooc", command, session.GetAccount().Username, session)
}
