package main

import (
	"akevitt/akevitt"
	"fmt"

	"github.com/rivo/tview"
)

func stats(engine *akevitt.Akevitt, session *ActiveSession) tview.Primitive {
	character := session.character

	format := fmt.Sprintf("HEALTH: %d/%d, NAME: %s (%s)", character.Health, character.MaxHealth, character.Name, character.currentRoom.GetName())
	return tview.NewTextView().SetText(format)
}

func visibleObjects(engine *akevitt.Akevitt, session *ActiveSession) tview.Primitive {
	return tview.NewList().
		AddItem("AAAA", "LOOOL", 0, nil).
		AddItem("AAAA", "LOOOL", 0, nil).
		AddItem("AAAA", "LOOOL", 0, nil).
		AddItem("AAAA", "LOOOL", 0, nil)
}
