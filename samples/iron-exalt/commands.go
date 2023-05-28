package main

import (
	"akevitt/akevitt"
	"errors"
	"fmt"
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

func characterStats(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {
	character, ok := session.RelatedGameObjects[currentCharacterKey].(*Character)
	if !ok {
		return errors.New("could not cast to character")
	}

	statsString := fmt.Sprintf(
		`
	===  %s (%s) ===
	  Health: %d/%d
	  Room: %s
	================
	`,
		character.Name,
		character.account.Username,
		character.Health,
		character.MaxHealth,
		character.currentRoom.Name)
	sepLines := strings.Split(statsString, "\n")
	for i := len(sepLines) - 1; i > 0; i-- {
		sender := "STATS"

		if i != 0 {
			sender = ""
		}

		AppendText(*session, sender, sepLines[i], ' ')
	}
	return nil
}

func help(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) error {

	helpString := `
	ooc <message> - Out-of-character Chat
	say <message> - Tell the message to everyone within the same room.
	stats		  - Show character stats
	help		  - Print this
	`

	sepLines := strings.Split(helpString, "\n")
	for i := len(sepLines) - 1; i > 0; i-- {
		sender := "STATS"

		if i != 0 {
			sender = ""
		}

		AppendText(*session, sender, sepLines[i], ' ')
	}
	return nil
}
