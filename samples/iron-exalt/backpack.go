package main

import (
	"akevitt/akevitt"
	"errors"
	"fmt"
	"strings"
)

func backpack(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*ActiveSession)

	if !ok {
		return errors.New("could not cast to session")
	}

	AppendText(sess, "Your backpack", sess.chat)
	for k, v := range sess.character.Inventory {
		AppendText(sess, fmt.Sprintf("№%d %s\n\t%s", k, v.GetName(), v.GetDescription()), sess.chat)
	}
	AppendText(sess, strings.Repeat("=.=", 16), sess.chat)

	return nil
}
