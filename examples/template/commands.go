package main

import (
	"akevitt"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Enter the room command
func enter(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*TemplateSession)
	if !ok {
		return errors.New("invalid session type")
	}
	character := sess.character
	prevRoom := character.currentRoom.GetName()
	roomKey, err := strconv.ParseUint(arguments, 10, 64)
	if err != nil {
		return err
	}
	exit, err := akevitt.IsRoomReachable[*Room](engine, session, roomKey, character.CurrentRoomKey)
	if err != nil {
		return err
	}
	err = exit.Enter(engine, session)

	if err != nil {
		return err
	}
	engine.Message(prevRoom, "left room", character.Name, session)
	engine.Message(character.currentRoom.GetName(), "entered room", character.Name, session)
	return nil
}

// Standard look command
func look(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*TemplateSession)

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

// Interact with an NPC or any other interactable objects
func interact(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*TemplateSession)

	if !ok {
		return errors.New("invalid session type")
	}

	arguments = strings.TrimSpace(arguments)
	for _, v := range akevitt.LookupOfType[Interactable](sess.character.currentRoom) {
		if !strings.EqualFold(v.GetName(), arguments) {
			continue
		}

		return v.Interact(engine, sess)
	}

	return fmt.Errorf("the object %s not found", arguments)
}

// Say command
func say(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*TemplateSession)

	if !ok {
		return errors.New("invalid session type")
	}

	return engine.Message(sess.character.currentRoom.GetName(), arguments, sess.character.Name, session)
}

// Out-of-character chat command
func ooc(engine *akevitt.Akevitt, session akevitt.ActiveSession, command string) error {
	return engine.Message("ooc", command, session.GetAccount().Username, session)
}

// View inventory
func backpack(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*TemplateSession)

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

// Mine ores
func mine(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*TemplateSession)

	if !ok {
		return errors.New("could not cast to session")
	}

	for _, v := range akevitt.LookupOfType[*Ore](sess.character.currentRoom) {
		if !strings.EqualFold(v.GetName(), arguments) {
			continue
		}

		return v.Use(engine, sess, sess.character.Inventory[0])
	}

	return nil
}
