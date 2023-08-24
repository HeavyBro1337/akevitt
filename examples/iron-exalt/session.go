package main

import (
	"github.com/IvanKorchmit/akevitt"

	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

// The custom session which holds reference to some UI primitives, subscribed channels of current user, character
// and account
type IronExaltSession struct {
	account            *akevitt.Account
	app                *tview.Application
	previousUI         *tview.Primitive
	subscribedChannels []string
	chat               *logview.LogView
	input              *tview.InputField
	character          *Character
	proceed            chan struct{}
}

func (sess *IronExaltSession) GetAccount() *akevitt.Account {
	return sess.account
}

func (sess *IronExaltSession) SetAccount(acc *akevitt.Account) {
	sess.account = acc
}

func (sess *IronExaltSession) GetApplication() *tview.Application {
	return sess.app
}

func (sess *IronExaltSession) SetApplication(app *tview.Application) {
	sess.app = app
}

func (sess *IronExaltSession) SetRoot(p tview.Primitive) {
	sess.previousUI = &p
	sess.app.SetRoot(p, true)
}
