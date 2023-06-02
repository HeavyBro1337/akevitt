package main

// This is the sample game using Akevitt
// Written by Ivan Korchmit (c) 2023

import (
	"akevitt/akevitt"
	"fmt"
	"log"
)

func main() {
	room := &Room{Name: "Spawn Room", DescriptionData: "Just a spawn room.", Key: 0, Exits: []uint64{1, 2, 3}}
	rooms := []*Room{
		{Name: "Mine", DescriptionData: "Mine of the corporation.", Key: 1, Exits: []uint64{1, 2, 3}},
		{Name: "Iron City", DescriptionData: "The lounge of the miners.", Key: 2, Exits: []uint64{1, 2, 3}},
	}

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
		RegisterCommand("enter", enterRoom).
		SetSpawnRoom(room)

	events := akevitt.GameEventHandler{}

	events.
		OOCMessage(func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, sender *akevitt.ActiveSession, message string) {
			err := AppendText(*session, fmt.Sprintf("%s (OOC): %s", sender.Account.Username, message))
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

			err := AppendText(*session, fmt.Sprintf("%s (Room): %s", senderCharacter.CharacterName, message))
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

		}).
		OnDatabaseCreate(func(engine *akevitt.Akevitt) error {
			fmt.Println("Database didn't exist. Creating rooms...")
			for _, v := range rooms {
				err := v.Save(v.Key, engine)
				if err != nil {
					return err
				}
			}
			return nil
		}).
		Finish()

	engine.ConfigureCallbacks(&events)

	log.Fatal(engine.Run())
}
