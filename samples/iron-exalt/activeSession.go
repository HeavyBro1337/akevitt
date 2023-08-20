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
	lookList           *tview.List
	character          *Character
}

func (sess *ActiveSession) GetAccount() *akevitt.Account {
	return sess.account
}

func (sess *ActiveSession) SetAccount(acc *akevitt.Account) {
	sess.account = acc
}

func (sess *ActiveSession) AppendLook(gameObject akevitt.GameObject) {
	if sess.lookList == nil {
		return
	}
	sess.lookList.AddItem(gameObject.GetName(), gameObject.GetDescription(), 0, nil)
}

func (sess *ActiveSession) ClearLook() {
	if sess.lookList == nil {
		return
	}
	sess.lookList.Clear()
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
