package main

import (
	"akevitt/akevitt"

	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

type ActiveSession struct {
	account            *akevitt.Account
	app                *tview.Application
	previousUI         *tview.Primitive
	subscribedChannels []string
	chat               *logview.LogView
	input              *tview.InputField
	character          *Character
	proceed            chan struct{}
}

func (sess *ActiveSession) GetAccount() *akevitt.Account {
	return sess.account
}

func (sess *ActiveSession) SetAccount(acc *akevitt.Account) {
	sess.account = acc
}

func (sess *ActiveSession) GetApplication() *tview.Application {
	return sess.app
}

func (sess *ActiveSession) SetApplication(app *tview.Application) {
	sess.app = app
}

func (sess *ActiveSession) GetPreviousUI() *tview.Primitive {
	return sess.previousUI
}

func (sess *ActiveSession) SetRoot(p tview.Primitive) {
	sess.previousUI = &p
	sess.app.SetRoot(p, true)
}
