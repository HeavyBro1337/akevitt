package basic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IvanKorchmit/akevitt"
)

// Enter the room command
func EnterCmd(engine *akevitt.Akevitt, session *Session, arguments string) error {
	character := session.Character
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

// Standard LookCmd command
func LookCmd(engine *akevitt.Akevitt, session *Session, arguments string) error {

	if strings.TrimSpace(arguments) == "" {
		for _, v := range session.Character.currentRoom.GetObjects() {
			AppendText(session, fmt.Sprintf("%s\n\t%s\n", v.GetName(), v.GetDescription()), session.Chat)
		}

		return nil
	}

	for _, v := range session.Character.currentRoom.GetObjects() {
		if strings.EqualFold(v.GetName(), arguments) {
			AppendText(session, fmt.Sprintf("%s\n\t%s\n", v.GetName(), v.GetDescription()), session.Chat)
		}
	}
	return nil
}

// Interact with an NPC or any other interactable objects
func TalkCmd(engine *akevitt.Akevitt, session *Session, arguments string) error {
	arguments = strings.TrimSpace(arguments)
	for _, v := range akevitt.LookupOfType[Interactable](session.Character.currentRoom) {
		if !strings.EqualFold(v.GetName(), arguments) {
			continue
		}

		return v.Interact(engine, session)
	}

	return fmt.Errorf("the object %s not found", arguments)
}

// Say command
func SayCmd(engine *akevitt.Akevitt, session *Session, arguments string) error {
	return engine.Message(session.Character.currentRoom.GetName(), arguments, session.Character.Name, session)
}

// Out-of-character chat command
func OocCmd(engine *akevitt.Akevitt, session *Session, command string) error {
	return engine.Message("ooc", command, session.GetAccount().Username, session)
}

// View inventory
func InventoryCmd(engine *akevitt.Akevitt, session *Session, arguments string) error {
	AppendText(session, "Your backpack", session.Chat)
	for k, v := range session.Character.Inventory {
		AppendText(session, fmt.Sprintf("№%d %s\n\t%s", k, v.GetName(), v.GetDescription()), session.Chat)
	}
	AppendText(session, strings.Repeat("=.=", 16), session.Chat)

	return nil
}
