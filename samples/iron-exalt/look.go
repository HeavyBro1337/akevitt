package main

import (
	"akevitt/akevitt"
	"errors"
	"fmt"
	"strings"
)

func look(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*ActiveSession)

	if !ok {
		return errors.New("could not cast to session")
	}
	if strings.TrimSpace(arguments) == "" {
		for _, v := range engine.Lookup(sess.character.currentRoom) {
			AppendText(sess, fmt.Sprintf("%s\n\t%s\n", v.GetName(), v.GetDescription()), sess.chat)
		}

		return nil
	}

	for _, v := range engine.Lookup(sess.character.currentRoom) {
		if strings.EqualFold(v.GetName(), arguments) {
			AppendText(sess, fmt.Sprintf("%s\n\t%s\n", v.GetName(), v.GetDescription()), sess.chat)
		}
	}
	return nil
}
