package main

// This is the sample game using Akevitt
// Written by Ivan Korchmit (c) 2023

import (
	"akevitt"
	"akevitt/core/network"
	"fmt"
	"log"

	"github.com/rivo/tview"
)

func main() {
	engine := akevitt.Akevitt{}
	engine.Defaults().
		GameName("Iron Exalt").
		Handle(nil).
		UseMouse(true).
		RootScreen(rootScreen).
		DatabasePath("data/iron-exalt.db").
		CreateDatabaseIfNotExists()

	events := akevitt.GameEventHandler{}

	events.
		OOCMessage(func(engine *akevitt.Akevitt, session *network.ActiveSession, sender *network.ActiveSession, message string) {
			err := AppendText(*session, *sender, message)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return
			}
			log.Println(sender.Account.Username)
			log.Println(message)
		}).
		Finish()

	engine.ConfigureCallbacks(&events)

	log.Fatal(engine.Run())
}

func rootScreen(engine *akevitt.Akevitt, session *network.ActiveSession) tview.Primitive {
	welcome := tview.NewModal().
		SetText(fmt.Sprintf("Welcome to %s. Would you like to register an account?", engine.GetGameName())).
		AddButtons([]string{"Register", "Login"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Login" {
				session.SetRoot(loginScreen(engine, session))
			} else if buttonLabel == "Register" {
				session.SetRoot(registerScreen(engine, session))
			}
		})
	return welcome
}
