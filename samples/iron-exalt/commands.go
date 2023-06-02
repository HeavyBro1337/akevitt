package main

import (
	"akevitt/akevitt"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ooc(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {
	engine.SendOOCMessage(command, session)

	return nil
}

func characterMessage(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {
	engine.SendRoomMessage(command, session)

	return nil
}

func enterRoom(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {
	character, ok := session.RelatedGameObjects[currentCharacterKey].Second.(*Character)
	if !ok {
		return errors.New("could not cast to character")
	}
	roomKey, err := strconv.ParseUint(command, 10, 64)

	if err != nil {
		return err
	}

	if engine.IsRoomReachable(roomKey, character.CurrentRoomKey) {
		character.CurrentRoomKey = roomKey
	}

	return nil
}

func characterStats(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {
	character, ok := session.RelatedGameObjects[currentCharacterKey].Second.(*Character)
	if !ok {
		return errors.New("could not cast to character")
	}

	statsString := fmt.Sprintf(
		`===  %s (%s) ===
  		  Health: %d/%d
  			Room: %s
		================`,
		character.CharacterName,
		character.account.Username,
		character.Health,
		character.MaxHealth,
		character.currentRoom.Name)
	sepLines := strings.Split(statsString, "\n")

	for _, v := range sepLines {

		AppendText(*session, v)
	}
	return nil
}

func help(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {

	helpString := `ooc <message> - Out-of-character Chat
	say <message> - Tell the message to everyone within the same room.
	look		  - Look around all the objects
	stats		  - Show character stats
	help		  - Print this`

	sepLines := strings.Split(helpString, "\n")
	for _, v := range sepLines {
		AppendText(*session, v)
	}
	return nil
}

func look(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {
	character, ok := session.RelatedGameObjects[currentCharacterKey].Second.(*Character)
	if !ok {
		return errors.New("could not cast to character")
	}

	// args := strings.Fields(command)

	if true {
		keysAndGo := akevitt.Lookup[*Character](engine, character.CurrentRoomKey)

		for _, v := range keysAndGo {
			err := AppendText(*session, fmt.Sprintf("%d -> %s", v.First, v.Second.Name()))

			if err != nil {
				return err
			}
		}
	}

	return nil
}
