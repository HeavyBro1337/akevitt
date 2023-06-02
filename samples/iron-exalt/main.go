package main

// This is the sample game using Akevitt
// Written by Ivan Korchmit (c) 2023

import (
	"akevitt/akevitt"
	"fmt"
	"log"
)

func main() {
	room := &Room{Name: "Spawn Room", DescriptionData: "Just a spawn room.", Key: 0}

	engine := akevitt.Akevitt{}
	engine.
		UseDefaults().
		UseGameName("Iron Exalt").
		UseMouse(true).
		UseRootScreen(rootScreen).
		UseDatabasePath("data/iron-exalt.db").
		UseCreateDatabaseIfNotExists().
		RegisterCommand("ooc", ooc).
		RegisterCommand("say", characterMessage).
		RegisterCommand("stats", characterStats).
		RegisterCommand("help", help).
		RegisterCommand("look", look).
		SetSpawnRoom(room)

	events := akevitt.GameEventHandler{}

	events.
		OOCMessage(func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, sender *akevitt.ActiveSession, message string) {
			err := AppendText(*session, sender.Account.Username, message, "%s (OOC): %s")
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
		}).
		Message(func(engine *akevitt.Akevitt, session, sender *akevitt.ActiveSession, message string) {
			senderCharacter, ok := sender.RelatedGameObjects[currentCharacterKey].Second.(*Character)
			if !ok {
				fmt.Println("Error: the current character turned out not to be a character struct!")
				return
			}

			sessionChacter, ok := session.RelatedGameObjects[currentCharacterKey].Second.(*Character)
			if !ok {
				fmt.Println("Error: the current character turned out not to be a character struct!")
				return
			}

			if sessionChacter.currentRoom.Name != senderCharacter.currentRoom.Name {
				return
			}

			err := AppendText(*session, senderCharacter.CharacterName, message, "%s (Room): %s")
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

		}).
		Finish()

	engine.ConfigureCallbacks(&events)

	log.Fatal(engine.Run())
}
