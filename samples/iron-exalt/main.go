package main

// This is the sample game using Akevitt
// Written by Ivan Korchmit (c) 2023

import (
	"akevitt"
	"akevitt/core/input"
	"akevitt/core/network"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	engine := akevitt.Akevitt{}
	engine.Defaults().
		GameName("Iron Exalt").
		Handle(nil).
		UseMouse().
		RootScreen(rootScreen).
		DatabasePath("data/iron-exalt.db").
		CreateDatabaseIfNotExists()

	engine.ConfigureCallbacks().
		LoginFail(func(engine *akevitt.Akevitt, session network.ActiveSession) {
			session.UI.SetRoot(ErrorBox("Login failed!!!!", session.UI, session.UIPrimitive), true)
		}).
		OOCMessage(func(engine *akevitt.Akevitt, session network.ActiveSession, sender network.ActiveSession, message string) {
			AppendText(session, sender, message)
		})

	log.Fatal(engine.Run())
}

func rootScreen(engine *akevitt.Akevitt, session network.ActiveSession) *tview.Modal {
	welcome := tview.NewModal().
		SetText(fmt.Sprintf("Welcome to %s. Would you like to register an account?", engine.GetGameName())).
		AddButtons([]string{"Yes", "Login"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Login" {
				session.UI.SetRoot(loginScreen(engine, session), true)
			} else if buttonLabel == "Yes" {
			}
		})
	return welcome
}

func loginScreen(engine *akevitt.Akevitt, session network.ActiveSession) *tview.Form {
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
		engine.Login(username, password, session)
	})
	return loginScreen
}

func ErrorBox[T tview.Primitive](message string, app *tview.Application, back *T) *tview.Modal {

	result := tview.NewModal().SetText("Error!").SetText(message).SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.SetRoot(*back, false)
	})
	result.SetBorderColor(tcell.ColorDarkRed)
	result.SetBorder(true)
	return result
}

func gameScreen(engine *akevitt.Akevitt, session network.ActiveSession) *tview.Grid {
	var playerMessage string
	const LABEL string = "Message: "
	chatLog := tview.NewList()
	inputField := tview.NewForm().AddInputField(LABEL, "", 32, nil, func(text string) {
		playerMessage = text
	})
	gameScreen := tview.NewGrid().
		SetRows(3).
		SetColumns(30).
		AddItem(chatLog, 1, 0, 3, 3, 0, 0, false).
		SetBorders(true).
		AddItem(inputField, 0, 0, 1, 3, 0, 0, true)
	inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if playerMessage == "" {
				return
			}
			switch status, parsedInput := input.ParseInput(playerMessage); status {
			case input.Message:
				if strings.TrimSpace(parsedInput) == "" {
					break
				}
				engine.SendOOCMessage(playerMessage, session)
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

func AppendText(currentSession network.ActiveSession, senderSession network.ActiveSession, message string) error {
	if currentSession.Chat == nil {
		return errors.New("TextView is nil")
	}
	currentSession.Chat.InsertItem(0, senderSession.Account.Username, message, 'M', nil)
	currentSession.Chat.SetWrapAround(true)
	return nil
}
