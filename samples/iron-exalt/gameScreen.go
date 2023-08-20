package main

import (
	"akevitt/akevitt"
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

func gameScreen(engine *akevitt.Akevitt, session *ActiveSession) tview.Primitive {

	var playerMessage string
	chatlog := logview.NewLogView()
	chatlog.SetLevelHighlighting(true)
	session.subscribedChannels = append(session.subscribedChannels, "ooc")
	session.chat = chatlog
	inputField := tview.NewInputField().SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}
		for _, word := range engine.GetCommands() {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
				entries = append(entries, word)
			}
		}

		f, ok := autocompletion[strings.Split(currentText, " ")[0]]

		if ok {
			for _, word := range f(currentText, engine, session) {
				if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
					entries = append(entries, word)
				}
			}
		}
		return entries
	}).SetChangedFunc(func(text string) {
		playerMessage = text
	})

	status := stats(engine, session)
	visibles := visibleObjects(engine, session)
	session.app.SetAfterDrawFunc(func(screen tcell.Screen) {
		lookupUpdate(engine, session, &visibles)
		fmt.Fprint(status.Clear(), updateStats(engine, session))
	})
	gameScreen := tview.NewGrid().
		SetRows(3).
		SetColumns(30).
		SetBorders(true).
		AddItem(inputField, 0, 0, 1, 1, 0, 0, true).
		AddItem(visibles, 1, 0, 1, 1, 0, 0, false).
		AddItem(status, 0, 1, 1, 2, 0, 0, false).
		AddItem(chatlog, 1, 1, 1, 2, 0, 0, false)

	inputField.SetFinishedFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			playerMessage = strings.TrimSpace(playerMessage)
			if playerMessage == "" {
				return
			}
			AppendText(session, playerMessage, session.chat)
			err := engine.ProcessCommand(playerMessage, session)
			if err != nil {
				ErrorBox(err.Error(), session.app, session.GetPreviousUI())
				inputField.Blur()
				inputField.SetText("")
				return
			}
			playerMessage = ""
			inputField.SetText("")
			session.app.SetFocus(inputField)
			lookupUpdate(engine, session, &visibles)
			fmt.Fprint(status.Clear(), updateStats(engine, session))
		}
	})
	inputField.SetAutocompletedFunc(func(text string, index, source int) bool {
		if source != tview.AutocompletedNavigate {
			inputField.SetText(text)
		}
		return source == tview.AutocompletedEnter || source == tview.AutocompletedClick
	})

	return gameScreen
}
