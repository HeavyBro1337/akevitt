package main

// This is the sample game using Akevitt
// Written by Ivan Korchmit (c) 2023

import (
	"akevitt/akevitt"
	"encoding/gob"
	"fmt"
	"log"
)

func main() {
	gob.Register(&Exit{})
	gob.Register(&Room{})

	room := &Room{Name: "Spawn Room", DescriptionData: "Just a spawn room.", Key: 0}
	rooms := []akevitt.Room{
		room,
		&Room{Name: "Mine", DescriptionData: "Mine of the corporation.", Key: 1},
		&Room{Name: "Iron City", DescriptionData: "The lounge of the miners.", Key: 2},
	}
	emptyExit := Exit{}
	BindRooms[*Exit](room, &emptyExit, rooms...)
	BindRooms[*Exit](rooms[1], &emptyExit, rooms...)
	BindRooms[*Exit](rooms[2], &emptyExit, rooms...)

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
		RegisterCommand("help", help).
		RegisterCommand("look", look).
		RegisterCommand("enter", enterRoom).
		SetSpawnRoom(room).
		SetStaticObjects(func(static map[uint64]akevitt.Object) {
			for _, v := range rooms {
				static[v.GetKey()] = v
			}
		})

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

			if sessionChacter.currentRoom.GetKey() != senderCharacter.currentRoom.GetKey() {
				return
			}

			err := AppendText(*session, fmt.Sprintf("%s (Room): %s", senderCharacter.CharacterName, message))
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

		}).
		OnDatabaseCreate(func(engine *akevitt.Akevitt) error {
			return nil
		}).
		Finish()

	engine.ConfigureCallbacks(&events)

	log.Fatal(engine.Run())
}
