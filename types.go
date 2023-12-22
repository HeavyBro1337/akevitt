package akevitt

import (
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

type pair[TFirst any, TSecond any] struct {
	f TFirst
	s TSecond
}

type Sessions = map[ssh.Session]*ActiveSession

type UIFunc = func(engine *Akevitt, session *ActiveSession) tview.Primitive

type CommandFunc = func(engine *Akevitt, session *ActiveSession, arguments string) error

type MessageFunc = func(engine *Akevitt, session *ActiveSession, channel, message, username string) error

type DeadSessionFunc = func(deadSession *ActiveSession, liveSessions []*ActiveSession, engine *Akevitt)

type DialogueFunc = func(engine *Akevitt, session *ActiveSession, dialogue *Dialogue) error
