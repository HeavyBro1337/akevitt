package basic

import (
	"github.com/IvanKorchmit/akevitt"
	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

const DEFAULT_KEY string = "CUSTOM_DATA"

type DefaultSessionData struct {
	akevitt.ActiveSession
	subscribedChannels []string
	account            *akevitt.Account
	previousUI         *tview.Primitive
	Chat               *logview.LogView
	Input              *tview.InputField
	proceed            chan struct{}
}

func InitSession(session *akevitt.ActiveSession) {
	defSession := DefaultSessionData{
		subscribedChannels: []string{"ooc"},
		Chat:               logview.NewLogView(),
	}
	session.Data[DEFAULT_KEY] = &defSession
}
