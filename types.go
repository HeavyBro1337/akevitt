package akevitt

import (
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type Pair[TFirst any, TSecond any] struct {
	L TFirst
	R TSecond
}

type Sessions = map[ssh.Session]*ActiveSession

type UIFunc = func(engine *Akevitt, session *ActiveSession) tview.Primitive

type CommandFunc = func(engine *Akevitt, session *ActiveSession, arguments string) error

type DeadSessionFunc = func(deadSession *ActiveSession, liveSessions []*ActiveSession, engine *Akevitt)

type DialogueFunc = func(engine *Akevitt, session *ActiveSession, dialogue *Dialogue) error
