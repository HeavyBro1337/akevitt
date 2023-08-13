package main

import (
	"akevitt/akevitt"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func characterCreationWizard(engine *akevitt.Akevitt, session *ActiveSession) tview.Primitive {
	var name string
	characterCreator := tview.NewForm().AddInputField("Character Name: ", "", 32, nil, func(text string) {
		name = text
	})
	characterCreator.AddButton("Done", func() {
		if strings.TrimSpace(name) == "" {
			ErrorBox("character name must not be empty!", session.app, session.GetPreviousUI())
			return
		}
		characterParams := CharacterParams{}
		characterParams.name = name
		emptychar := &Character{}
		_, err := akevitt.CreateObject(engine, session, emptychar, characterParams)
		if err != nil {
			ErrorBox(err.Error(), session.app, session.GetPreviousUI())
			return
		}
		session.SetRoot(gameScreen(engine, session))
	})
	return characterCreator
}

func gameScreen(engine *akevitt.Akevitt, session *ActiveSession) tview.Primitive {
	fmt.Printf("session: %v\n", session.account)
	var playerMessage string
	const LABEL string = "Message: "
	chatlog := tview.NewList()

	session.chat = chatlog

	inputField := tview.NewForm().AddInputField(LABEL, "", 32, nil, func(text string) {
		playerMessage = text
	})
	gameScreen := tview.NewGrid().
		SetRows(3).
		SetColumns(30).
		AddItem(chatlog, 1, 0, 3, 3, 0, 0, false).
		SetBorders(true).
		AddItem(inputField, 0, 0, 1, 1, 0, 0, true)
		// AddItem(stats(engine, session), 0, 1, 1, 2, 0, 0, false)
	inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			playerMessage = strings.TrimSpace(playerMessage)
			if playerMessage == "" {
				return
			}
			err := AppendText(session, playerMessage)
			if err != nil {
				ErrorBox(err.Error(), session.app, session.GetPreviousUI())
				playerMessage = ""
				inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetText("")
				session.app.SetFocus(inputField.GetFormItemByLabel(LABEL))
				return
			}
			err = engine.ProcessCommand(playerMessage, session)
			if err != nil {
				ErrorBox(err.Error(), session.app, session.GetPreviousUI())
			}
			playerMessage = ""
			inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetText("")
			session.app.SetFocus(inputField.GetFormItemByLabel(LABEL))
		}
	})
	return gameScreen
}
