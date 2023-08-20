package main

import (
	"akevitt/akevitt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func DialogueBox(dial *akevitt.Dialogue, engine *akevitt.Akevitt, session *ActiveSession) error {
	labels := akevitt.MapSlice(dial.GetOptions(), func(v *akevitt.Dialogue) string {
		return v.GetTitle()
	})

	if len(labels) == 0 {
		session.app.SetRoot(*session.previousUI, true)
		session.proceed <- struct{}{}
		return nil
	}
	var err error = nil
	modal := tview.NewModal().AddButtons(labels).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		err = dial.Proceed(buttonIndex, session, engine)
		if len(dial.GetOptions()) == 0 {
			session.app.SetRoot(*session.previousUI, true)
		}
	})

	modal.SetBackgroundColor(tcell.ColorBlack).SetBorder(false)

	grid := func(p tview.Primitive) tview.Primitive {
		gr := tview.NewGrid().
			SetColumns(3).
			SetRows(3).
			AddItem(p, 0, 1, 3, 2, 0, 0, false).
			AddItem(modal, 1, 1, 1, 1, 0, 0, true)
		gr.SetBorder(true)
		return gr
	}

	session.app.SetRoot(grid(dial.GetContents()), true)
	return err
}
