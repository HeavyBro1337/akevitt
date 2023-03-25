package main

// This is the sample game using Akevitt
// Written by Ivan Korchmit (c) 2023

import (
	"akevitt/akevitt"
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/rivo/tview"
)

func main() {

	engine := akevitt.Akevitt{}
	engine.Defaults().
		GameName("Iron Exalt").
		UseMouse(true).
		RootScreen(rootScreen).
		DatabasePath("data/iron-exalt.db").
		CreateDatabaseIfNotExists().
		RegisterCommand("ooc", ooc)

	events := akevitt.GameEventHandler{}

	events.
		OOCMessage(func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, sender *akevitt.ActiveSession, message string) {
			err := AppendText(*session, *sender, message)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return
			}
		}).
		Finish()

	engine.ConfigureCallbacks(&events)

	log.Fatal(engine.Run())
}

func rootScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	b, err := os.ReadFile("./data/logo.png")
	if err != nil {
		panic("Cannot find image!!!")
	}
	pngLogo, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err.Error())
	}
	image := tview.NewImage().SetImage(pngLogo)
	wizard := tview.NewModal().
		SetText(fmt.Sprintf("Welcome to %s! Would you register your account?", engine.GetGameName())).
		AddButtons([]string{"Register", "Login"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Login" {
				session.SetRoot(loginScreen(engine, session))
			} else if buttonLabel == "Register" {
				session.SetRoot(registerScreen(engine, session))
			}
		})
	welcome := tview.NewGrid().
		SetBorders(false).
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		AddItem(image, 0, 0, 3, 27, 0, 0, false).
		AddItem(wizard, 2, 2, 3, 3, 0, 0, false)
	return welcome
}

type Room struct {
	Name string
}
type Character struct {
	Name      string
	Health    int
	MaxHealth int
	account   *akevitt.Account
}
