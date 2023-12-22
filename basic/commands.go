package basic

import (
	"fmt"
	"strings"

	"github.com/IvanKorchmit/akevitt"
)

// Enter the room command
func EnterCmd(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess := CastSession[*Session](session)

	character := sess.Character
	prevRoom := character.currentRoom.GetName()

	exit, err := akevitt.IsRoomReachable[*Room](engine, session, arguments, character.CurrentRoomKey)
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

// Standard LookCmd command
func LookCmd(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess := CastSession[*Session](session)

	if strings.TrimSpace(arguments) == "" {
		for _, v := range sess.Character.currentRoom.GetObjects() {
			AppendText(fmt.Sprintf("%s\n\t%s\n", v.GetName(), v.GetDescription()), sess.Chat)
		}

		return nil
	}

	for _, v := range sess.Character.currentRoom.GetObjects() {
		if strings.EqualFold(v.GetName(), arguments) {
			AppendText(fmt.Sprintf("%s\n\t%s\n", v.GetName(), v.GetDescription()), sess.Chat)
		}
	}
	return nil
}

// Interact with an NPC or any other interactable objects
func TalkCmd(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess := CastSession[*Session](session)

	arguments = strings.TrimSpace(arguments)
	for _, v := range akevitt.LookupOfType[Interactable](sess.Character.currentRoom) {
		if !strings.EqualFold(v.GetName(), arguments) {
			continue
		}

		return v.Interact(engine, sess)
	}

	return fmt.Errorf("the object %s not found", arguments)
}

// Say command
func SayCmd(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess := CastSession[*Session](session)

	return engine.Message(sess.Character.currentRoom.GetName(), arguments, sess.Character.Name, session)
}

// Out-of-character chat command
func OocCmd(engine *akevitt.Akevitt, session akevitt.ActiveSession, command string) error {
	return engine.Message("ooc", command, session.GetAccount().Username, session)
}

// View inventory
func InventoryCmd(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess := CastSession[*Session](session)

	AppendText("Your backpack", sess.Chat)
	for k, v := range sess.Character.Inventory {
		AppendText(fmt.Sprintf("№%d %s\n\t%s", k, v.GetName(), v.GetDescription()), sess.Chat)
	}
	AppendText(strings.Repeat("=.=", 16), sess.Chat)

	return nil
}
