package engine

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ActiveSession struct {
	Account     *Account
	Application *tview.Application
	Data        map[string]any
	RoomID      string
}

func (session *ActiveSession) Send(message string) {
	fmt.Println(message)
}

func (session *ActiveSession) Sendf(format string, args ...any) {
	session.Send(fmt.Sprintf(format, args...))
}

func (session *ActiveSession) SendLines(lines ...string) {
	session.Send(strings.Join(lines, "\n"))
}

func PurgeDeadSessions(engine *Akevitt, callback ...DeadSessionFunc) {
	deadSessions := make([]*ActiveSession, 0)
	liveSessions := make([]*ActiveSession, 0)
	for k, v := range engine.sessions {
		_, err := io.WriteString(k, "\000")
		if err != nil {
			delete(engine.sessions, k)
			deadSessions = append(deadSessions, v)
			continue
		}
		liveSessions = append(liveSessions, v)
	}
	if callback != nil {
		for _, v := range deadSessions {
			for _, fn := range callback {
				fn(v, liveSessions, engine)
			}
		}
	}
}

func AppendText(message string, chatlog io.Writer) {
	fmt.Fprintln(chatlog, message)
}

func ErrorBox(message string, app *tview.Application, back tview.Primitive) {
	result := tview.NewModal().SetText("Error!").SetText(message).SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.SetRoot(back, true)
	})

	result.SetBorder(true).SetBorderColor(tcell.ColorDarkRed)
	app.SetRoot(result, true)
}

const format string = "[%s] %s: %s\n"

func LogInfo(message string) {
	fmt.Printf(format, time.Now(), "LOG", message)
}
func LogWarn(message string) {
	fmt.Printf(format, time.Now(), "WARN", message)
}
func LogError(message string) {
	fmt.Printf(format, time.Now(), "ERR", message)
}