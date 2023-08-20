package main

import (
	"akevitt/akevitt"
	"strings"

	"github.com/rivo/tview"
)

func characterCreationWizard(engine *akevitt.Akevitt, session *ActiveSession) tview.Primitive {
	var name string
	characterCreator := tview.NewForm().AddInputField("Character Name: ", "", 32, nil, func(text string) {
		name = text
	})
	characterCreator.AddButton("Done", func() {
		if strings.TrimSpace(name) == "" {
			ErrorBox("character name must not be empty!", session, session.GetPreviousUI())
			return
		}
		characterParams := CharacterParams{}
		characterParams.name = name
		emptyChar := &Character{}
		_, err := akevitt.CreateObject(engine, session, emptyChar, characterParams)
		if err != nil {
			ErrorBox(err.Error(), session, session.GetPreviousUI())
			return
		}
		session.SetRoot(gameScreen(engine, session))
	})
	return characterCreator
}
