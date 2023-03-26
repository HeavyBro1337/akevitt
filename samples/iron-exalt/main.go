package main

// This is the sample game using Akevitt
// Written by Ivan Korchmit (c) 2023

import (
	"akevitt/akevitt"
	"fmt"
	"log"
)

func main() {

	engine := akevitt.Akevitt{}
	engine.
		UseDefaults().
		UseGameName("Iron Exalt").
		UseMouse(true).
		UseRootScreen(rootScreen).
		UseDatabasePath("data/iron-exalt.db").
		UseCreateDatabaseIfNotExists().
		RegisterCommand("ooc", ooc).
		RegisterCommand("say", characterMessage)

	events := akevitt.GameEventHandler{}

	events.
		OOCMessage(func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, sender *akevitt.ActiveSession, message string) {
			err := AppendText(*session, sender.Account.Username, message, 'O')
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}
		}).
		Message(func(engine *akevitt.Akevitt, session, sender *akevitt.ActiveSession, message string) {
			character, ok := sender.RelatedGameObjects[currentCharacterKey].(*Character)
			if !ok {
				fmt.Println("Error: the current character turned out not to be a character struct!")
				return
			}

			err := AppendText(*session, character.name, message, 'R')
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

		}).
		Finish()

	engine.ConfigureCallbacks(&events)

	log.Fatal(engine.Run())
}

type Room struct {
	Name string
}
