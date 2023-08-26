package basic

import (
	"github.com/IvanKorchmit/akevitt"

	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

// The custom session which holds reference to some UI primitives, subscribed channels of current user, character
// and account
type Session struct {
	account            *akevitt.Account
	app                *tview.Application
	previousUI         *tview.Primitive
	subscribedChannels []string
	Chat               *logview.LogView
	Input              *tview.InputField
	Character          *Character
	proceed            chan struct{}
}

func (sess *Session) GetAccount() *akevitt.Account {
	return sess.account
}

func (sess *Session) SetAccount(acc *akevitt.Account) {
	sess.account = acc
}

func (sess *Session) GetApplication() *tview.Application {
	return sess.app
}

func (sess *Session) SetApplication(app *tview.Application) {
	sess.app = app
}

func (sess *Session) SetRoot(p tview.Primitive) {
	sess.previousUI = &p
	sess.app.SetRoot(p, true)
}

func (sess *Session) GetCurrentUI() *tview.Primitive {
	return sess.previousUI
}

func CastSession[T akevitt.ActiveSession](session akevitt.ActiveSession) T {
	return session.(T)
}
