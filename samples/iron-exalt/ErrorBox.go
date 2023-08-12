package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func ErrorBox(message string, app *tview.Application, back *tview.Primitive) {
	result := tview.NewModal().SetText("Error!").SetText(message).SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.SetRoot(*back, false)
	}).SetFocus(0)
	result.SetBorderColor(tcell.ColorDarkRed)
	result.SetBorder(true)
	app.SetRoot(result, true)
}
