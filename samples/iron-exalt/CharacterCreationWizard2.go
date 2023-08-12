package main

import (
	"akevitt/akevitt"
	"strings"

	"github.com/rivo/tview"
)

func characterCreationWizard(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	var name string
	characterCreator := tview.NewForm().AddInputField("Character Name: ", "", 32, nil, func(text string) {
		name = text
	})
	characterCreator.AddButton("Done", func() {
		if strings.TrimSpace(name) == "" {
			ErrorBox("character name must not be empty!", session.UI, session.GetPreviousUI())
			return
		}
		characterParams := CharacterParams{}
		characterParams.name = name
		emptychar := &Character{}
		_, err := akevitt.CreateObject(engine, session, emptychar, characterParams)
		if err != nil {
			ErrorBox(err.Error(), session.UI, session.GetPreviousUI())
			return
		}
		session.SetRoot(gameScreen(engine, session))
	})
	return characterCreator
}
