package main

import (
	"akevitt/akevitt"
	"errors"
	"strings"
)

func interact(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*ActiveSession)

	if !ok {
		return errors.New("invalid session type")
	}

	return engine.Interact(strings.TrimSpace(arguments), sess.character.currentRoom, session)
}
