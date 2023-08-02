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

func lookExits(engine *akevitt.Akevitt, session *akevitt.ActiveSession) error {
	character, ok := session.RelatedGameObjects[currentCharacterKey].Second.(*Character)
	if !ok {
		return errors.New("could not cast to character")
	}
	for _, v := range character.currentRoom.GetExits() {
		room, err := akevitt.GetStaticObject[*Room](engine, v.GetKey())

		if err != nil {
			return err
		}

		err = AppendText(*session, fmt.Sprintf("Room %s (%d)", room.Name, v.GetKey()))

		if err != nil {
			return err
		}
	}
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

	exit, err := akevitt.IsRoomReachable[*Room](engine, session, roomKey, character.CurrentRoomKey)

	if err != nil {
		return err
	}

	return exit.Enter(engine, session)
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
		character.currentRoom.Description())
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

	key, err := strconv.ParseUint(command, 10, 64)

	if err != nil {
		keysAndGo := akevitt.Lookup[akevitt.GameObject](engine, character.CurrentRoomKey)

		err := lookExits(engine, session)

		if err != nil {
			return err
		}
		for _, v := range keysAndGo {
			err := AppendText(*session, fmt.Sprintf("%d -> %s", v.First, v.Second.Name()))

			if err != nil {
				return err
			}
		}
	} else {
		obj, err := akevitt.GetObject[akevitt.GameObject](engine, key, false)
		if err != nil {
			return err
		}

		if obj.OnRoomLookup() != character.CurrentRoomKey {
			return errors.New("unknown key specified")
		}
		return AppendText(*session, fmt.Sprintf("%s: %s", obj.Name(), obj.Description()))
	}

	return nil
}
