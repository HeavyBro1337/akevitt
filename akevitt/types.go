package akevitt

import (
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type Pair[TFirst any, TSecond any] struct {
	First  TFirst
	Second TSecond
}

type Sessions = map[ssh.Session]*ActiveSession

type UIFunc = func(engine *akevitt, session *ActiveSession) tview.Primitive

type CommandFunc = func(engine *akevitt, session *ActiveSession, command string) error
