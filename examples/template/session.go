package main

import (
	"akevitt"

	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

// The custom session which holds reference to some UI primitives, subscribed channels of current user, character
// and account
type TemplateSession struct {
	account            *akevitt.Account
	app                *tview.Application
	previousUI         *tview.Primitive
	subscribedChannels []string
	chat               *logview.LogView
	input              *tview.InputField
	character          *Character
	proceed            chan struct{}
}

func (sess *TemplateSession) GetAccount() *akevitt.Account {
	return sess.account
}

func (sess *TemplateSession) SetAccount(acc *akevitt.Account) {
	sess.account = acc
}

func (sess *TemplateSession) GetApplication() *tview.Application {
	return sess.app
}

func (sess *TemplateSession) SetApplication(app *tview.Application) {
	sess.app = app
}

func (sess *TemplateSession) SetRoot(p tview.Primitive) {
	sess.previousUI = &p
	sess.app.SetRoot(p, true)
}
