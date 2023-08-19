package main

import (
	"akevitt/akevitt"
	"fmt"

	"github.com/rivo/tview"
)

func stats(engine *akevitt.Akevitt, session *ActiveSession) *tview.TextView {
	return tview.NewTextView().SetText(updateStats(engine, session))
}

func updateStats(engine *akevitt.Akevitt, session *ActiveSession) string {
	character := session.character
	return fmt.Sprintf("HEALTH: %d/%d, NAME: %s (%s)", character.Health,
		character.MaxHealth,
		character.Name,
		character.currentRoom.GetName())
}

func visibleObjects(engine *akevitt.Akevitt, session *ActiveSession) tview.Primitive {
	return tview.NewList().
		AddItem("AAAA", "LOOOL", 0, nil).
		AddItem("AAAA", "LOOOL", 0, nil).
		AddItem("AAAA", "LOOOL", 0, nil).
		AddItem("AAAA", "LOOOL", 0, nil)
}
