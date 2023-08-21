package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ErrorBox(message string, session *ActiveSession, back *tview.Primitive) {
	result := tview.NewModal().SetText("Error!").SetText(message).SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		session.app.SetRoot(*back, true)
		if session.input != nil {
			session.app.SetFocus(session.input)
		}
	})

	result.SetBorder(true).SetBorderColor(tcell.ColorDarkRed)
	session.app.SetRoot(result, true)
}
