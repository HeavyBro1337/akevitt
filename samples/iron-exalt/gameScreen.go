package main

import (
	"akevitt/akevitt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func gameScreen(engine *akevitt.Akevitt, session *ActiveSession) tview.Primitive {
	var playerMessage string
	const LABEL string = "Message: "
	chatlog := tview.NewList()
	session.subscribedChannels = append(session.subscribedChannels, "ooc")
	session.chat = chatlog

	inputField := tview.NewForm().AddInputField(LABEL, "", 32, nil, func(text string) {
		playerMessage = text
	})
	gameScreen := tview.NewGrid().
		SetRows(3).
		SetColumns(30).
		AddItem(chatlog, 1, 0, 3, 3, 0, 0, false).
		SetBorders(true).
		AddItem(inputField, 0, 0, 1, 1, 0, 0, true).
		AddItem(stats(engine, session), 0, 1, 1, 2, 0, 0, false)
	inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			playerMessage = strings.TrimSpace(playerMessage)
			if playerMessage == "" {
				return
			}
			err := AppendText(session, playerMessage, session.chat)
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
				return
			}
			playerMessage = ""
			inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetText("")
			session.app.SetFocus(inputField.GetFormItemByLabel(LABEL))
		}
	})
	return gameScreen
}
