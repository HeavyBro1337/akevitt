package main

import (
	"akevitt"
	"akevitt/core/input"
	"akevitt/core/network"
	"errors"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func AppendText(currentSession network.ActiveSession, senderSession network.ActiveSession, message string) error {
	if currentSession.Chat == nil {
		return errors.New("Chat log element is nil")
	}
	currentSession.Chat.InsertItem(0, senderSession.Account.Username, message, 'M', nil)
	currentSession.Chat.SetWrapAround(true)
	return nil
}

func loginScreen(engine *akevitt.Akevitt, session *network.ActiveSession) tview.Primitive {
	var username string
	var password string
	loginScreen := tview.NewForm().
		AddInputField("Username: ", "", 32, nil, func(text string) {
			username = text
		}).
		AddPasswordField("Password: ", "", 32, '*', func(text string) {
			password = text
		})
	loginScreen.AddButton("Login", func() {
		err := engine.Login(username, password, session)

		if err != nil {
			fmt.Printf("err: %v\n", err)
			ErrorBox(err.Error(), session.UI, session.UIPrimitive)
			return
		}
		session.SetRoot(gameScreen(engine, session))
	})
	return loginScreen
}

func ErrorBox(message string, app *tview.Application, back *tview.Primitive) {
	fmt.Printf("app: %v\n", app)
	fmt.Printf("back: %v\n", back)
	result := tview.NewModal().SetText("Error!").SetText(message).SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.SetRoot(*back, false)
	})
	result.SetBorderColor(tcell.ColorDarkRed)
	result.SetBorder(true)
	app.SetRoot(result, true)
}

func gameScreen(engine *akevitt.Akevitt, session *network.ActiveSession) tview.Primitive {
	fmt.Printf("session: %v\n", session.Account)
	var playerMessage string
	const LABEL string = "Message: "
	session.Chat = tview.NewList()
	inputField := tview.NewForm().AddInputField(LABEL, "", 32, nil, func(text string) {
		playerMessage = text
	})
	gameScreen := tview.NewGrid().
		SetRows(3).
		SetColumns(30).
		AddItem(session.Chat, 1, 0, 3, 3, 0, 0, false).
		SetBorders(true).
		AddItem(inputField, 0, 0, 1, 3, 0, 0, true)
	inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if playerMessage == "" {
				return
			}
			status, parsedInput := input.ParseInput(playerMessage)
			switch status {
			case input.Message:
				if strings.TrimSpace(parsedInput) == "" {
					break
				}
				engine.SendOOCMessage(parsedInput, session)
			case input.Command:
				if strings.TrimSpace(parsedInput) == "" {
					break
				}
			}

			playerMessage = ""
			inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetText("")
			session.UI.SetFocus(inputField.GetFormItemByLabel(LABEL))
		}
	})
	return gameScreen
}

func registerScreen(engine *akevitt.Akevitt, session *network.ActiveSession) tview.Primitive {
	var username string
	var password string
	var repeatPassword string

	gameScreen := gameScreen(engine, session)
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

				ErrorBox("Passwords don't match!", session.UI, session.UIPrimitive)
				return
			}
			err := engine.Register(username, password, session)
			if err != nil {
				ErrorBox(err.Error(), session.UI, session.UIPrimitive)
				return
			}
			session.SetRoot(gameScreen)
		})
	registerScreen.SetBorder(true).SetTitle("Register")
	return registerScreen
}
