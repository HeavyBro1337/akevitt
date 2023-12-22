package akevitt

import (
	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

const DEFAULT_KEY string = "CUSTOM_DATA"

type DefaultSessionData struct {
	subscribedChannels []string
	previousUI         *tview.Primitive
	Chat               *logview.LogView
	Input              *tview.InputField
	proceed            chan struct{}
}

func InitSession(session *ActiveSession) {
	defSession := DefaultSessionData{
		subscribedChannels: []string{"ooc"},
		Chat:               logview.NewLogView(),
	}
	session.Data[DEFAULT_KEY] = &defSession
}
