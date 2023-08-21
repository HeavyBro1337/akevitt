package main

import (
	"akevitt/akevitt"
	"fmt"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func stats(engine *akevitt.Akevitt, session *ActiveSession) *tview.TextView {
	return tview.NewTextView().SetText(updateStats(engine, session))
}

func updateStats(engine *akevitt.Akevitt, session *ActiveSession) string {
	character := session.character
	return fmt.Sprintf("HEALTH: %d/%d, NAME: %s (%s) $%d", character.Health,
		character.MaxHealth,
		character.Name,
		character.currentRoom.GetName(),
		character.Money)
}

func visibleObjects(engine *akevitt.Akevitt, session *ActiveSession) *tview.List {
	l := tview.NewList()
	lookupUpdate(engine, session, &l)
	return l
}

func lookupUpdate(engine *akevitt.Akevitt, session *ActiveSession, l **tview.List) {
	objects := engine.Lookup(session.character.currentRoom)
	(*l).Clear()
	for _, v := range objects {
		if v == session.character {
			continue
		}

		(*l).AddItem(v.GetName(), v.GetDescription(), 0, nil)
	}
	(*l).AddItem("AVAILABLE ROOMS", "", 0, nil)
	exits := session.character.currentRoom.GetExits()

	for _, v := range exits {
		(*l).AddItem(v.GetRoom().GetName(), strconv.FormatUint(v.GetKey(), 10), 0, nil)
	}
	(*l).SetSelectedBackgroundColor(tcell.ColorBlack).SetSelectedTextColor(tcell.ColorWhite)
}
