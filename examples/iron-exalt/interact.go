package main

import (
	"akevitt/akevitt"
	"errors"
	"fmt"
	"strings"
)

func interact(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*ActiveSession)

	if !ok {
		return errors.New("invalid session type")
	}

	arguments = strings.TrimSpace(arguments)
	interactables := akevitt.LookupOfType[Interactable](sess.character.currentRoom)
	for _, v := range interactables {
		if !strings.EqualFold(v.GetName(), arguments) {
			continue
		}

		return v.Interact(engine, sess)
	}

	return fmt.Errorf("the object %s not found", arguments)
}
