package akevitt

import (
	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

const DEFAULT_KEY string = "CUSTOM_DATA"

type DefaultSessionData struct {
	SubscribedChannels []string
	previousUI         *tview.Primitive
	Chat               *logview.LogView
	Input              *tview.InputField
}

func InitSession(session *ActiveSession) {
	defSession := DefaultSessionData{
		SubscribedChannels: []string{"ooc"},
		Chat:               logview.NewLogView(),
		Input:              tview.NewInputField(),
	}
	session.Data[DEFAULT_KEY] = &defSession
}
