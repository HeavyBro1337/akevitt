package main

import (
	"akevitt/akevitt"

	"github.com/rivo/tview"
)

func registerScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	var username string
	var password string
	var repeatPassword string

	registerScreen := tview.NewForm().AddInputField("Username: ", "", 32, nil, func(text string) {
		username = text
	}).
		AddPasswordField("Password: ", "", 32, '*', func(text string) {
			password = text
		}).
		AddPasswordField("Repeat password: ", "", 32, '*', func(text string) {
			repeatPassword = text
		})
	registerScreen.
		AddButton("Create account", func() {
			if password != repeatPassword {
				ErrorBox("Passwords don't match!", session.UI, session.GetPreviousUI())
				return
			}
			err := engine.Register(username, password, session)
			if err != nil {
				ErrorBox(err.Error(), session.UI, session.GetPreviousUI())
				return
			}
			session.SetRoot(characterCreationWizard(engine, session))
		})
	registerScreen.SetBorder(true).SetTitle(" Register ")
	return registerScreen
}
