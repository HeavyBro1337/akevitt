package akevitt

import (
	"unicode"

	"github.com/IvanKorchmit/akevitt"
	"github.com/IvanKorchmit/akevitt/plugins"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

func AppendText(message string, chatlog *logview.LogView) {
	ev := logview.NewLogEvent("message", message)
	ev.Level = logview.LogLevelInfo
	chatlog.AppendEvent(ev)
	chatlog.SetFocusFunc(func() {
		chatlog.Blur()
	})
	chatlog.ScrollToBottom()
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

func RegistrationScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession, nextScreen akevitt.UIFunc) tview.Primitive {
	username := ""
	password := ""
	repeatPassword := ""

	account := akevitt.FetchPluginUnsafe[*plugins.AccountPlugin](engine)

	form := tview.NewForm()

	form.AddInputField("Username", "", 0, func(textToCheck string, lastChar rune) bool {
		if !unicode.IsLetter(lastChar) && !unicode.IsDigit(lastChar) || lastChar > unicode.MaxASCII {
			return false
		}

		username = textToCheck
		return true
	}, nil).
		AddPasswordField("Repeat password", "", 0, '*', func(text string) {
			password = text
		}).
		AddPasswordField("Repeat password", "", 0, '*', func(text string) {
			repeatPassword = text
		}).
		AddButton("Register", func() {
			err := account.Register(username, password, repeatPassword, session)

			if err != nil {
				ErrorBox(err.Error(), session.Application, form)
				return
			}

			session.Application.SetRoot(nextScreen(engine, session), true)
		})

	return form
}

func LoginScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession, nextScreen akevitt.UIFunc) tview.Primitive {
	username := ""
	password := ""

	form := tview.NewForm()

	account := akevitt.FetchPluginUnsafe[*plugins.AccountPlugin](engine)

	form.AddInputField("Username", "", 0, func(textToCheck string, lastChar rune) bool {
		if !unicode.IsLetter(lastChar) && !unicode.IsDigit(lastChar) || lastChar > unicode.MaxASCII {
			return false
		}

		username = textToCheck
		return true
	}, nil).
		AddPasswordField("Password", "", 0, '*', func(text string) {
			password = text
		}).
		AddButton("Login", func() {
			err := account.Login(username, password, session)

			if err != nil {
				ErrorBox(err.Error(), session.Application, form)
				return
			}

			session.Application.SetRoot(nextScreen(engine, session), true)
		})

	form.SetTitle("Login")

	return form
}
