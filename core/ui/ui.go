package ui

import (
	"akevitt/core/database"
	"akevitt/core/database/credentials"
	"akevitt/core/input"
	"akevitt/core/network"
	"errors"

	"github.com/boltdb/bolt"
	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

// generate register screen
func generateRegistrationScreen(sesh ssh.Session, sessions *map[ssh.Session]network.ActiveSession, db *bolt.DB) (*tview.Form, bool) {
	var username string
	var password string
	var repeatPassword string

	gameScreen, chatLog := GenerateGameScreen(sesh, sessions)

	registerScreen := tview.NewForm().AddInputField("Username: ", "", 32, nil, func(text string) {
		username = text
	}).
		AddPasswordField("Password: ", "", 32, '*', func(text string) {
			password = text
		}).
		AddPasswordField("Repeat password: ", "", 32, '*', func(text string) {
			repeatPassword = text
		})

	registerScreen.
		AddButton("Create account", func() {
			network.PurgeDeadSessions(sessions)
			if password != repeatPassword {
				(*sessions)[sesh].UI.SetRoot(ErrorBox("Password don't match", (*sessions)[sesh].UI, &registerScreen),
					false)
				return
			}
			if !database.DoesAccountExist(username, db) {
				id, err := database.CreateAccount(db, username, password)
				if err != nil {
					return
				}
				if active, ok := (*sessions)[sesh]; ok {
					acc, err := database.GetAccount(id, db)
					if err != nil {
						return
					}
					active.Account = &acc
					active.Chat = chatLog
					(*sessions)[sesh] = active
					(*sessions)[sesh].UI.SetRoot(gameScreen, true)
				}
			} else {
				(*sessions)[sesh].UI.SetRoot(ErrorBox("Account already exists", (*sessions)[sesh].UI, &registerScreen), false)
			}
		})
	registerScreen.SetBorder(true).SetTitle("Register")
	return registerScreen, true
}

// Appends message to the end-user
func AppendText(text *tview.List, account credentials.Account, message string, ui *tview.Application) error {
	if text == nil {
		return errors.New("TextView is nil")
	}
	text.InsertItem(0, account.Username, message, 'M', nil)
	text.SetWrapAround(true)
	return nil
}

func AppendNoSenderText(text *tview.List, message string, ui *tview.Application) error {
	if text == nil {
		return errors.New("TextView is nil")
	}
	text.
		AddItem("", message, 'S', nil)

	return nil
}

// generate login screen
func generateLoginScreen(sesh ssh.Session, sessions *map[ssh.Session]network.ActiveSession, db *bolt.DB) *tview.Form {
	var username string
	var password string

	gameScreen, chatLog := GenerateGameScreen(sesh, sessions)

	loginScreen := tview.NewForm().AddInputField("Username: ", "", 32, nil, func(text string) {
		username = text
	}).
		AddPasswordField("Password: ", "", 32, '*', func(text string) {
			password = text
		})

	loginScreen.
		AddButton("Login", func() {
			network.PurgeDeadSessions(sessions)
			ok, acc := database.Login(username, password, db)
			if ok {
				if !database.CheckCurrentLogin(*acc, sessions) {
					if active, ok := (*sessions)[sesh]; ok {
						active.Account = acc
						active.Chat = chatLog
						(*sessions)[sesh] = active
					}
					(*sessions)[sesh].UI.SetRoot(gameScreen, true)

				} else {
					(*sessions)[sesh].UI.SetRoot(ErrorBox("Somebody already logged in!", (*sessions)[sesh].UI, &loginScreen), false)
				}
			} else {
				(*sessions)[sesh].UI.SetRoot(ErrorBox("Wrong password or username!", (*sessions)[sesh].UI, &loginScreen), false)
			}

		})
	loginScreen.SetBorder(true).SetTitle("Login")
	return loginScreen
}
func ErrorBox[T tview.Primitive](message string, app *tview.Application, back *T) *tview.Modal {
	result := tview.NewModal().SetText("Error!").SetText(message).SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.SetRoot(*back, false)
	})
	result.SetBorderColor(tcell.ColorDarkRed)
	result.SetBorder(true)
	return result
}

// generate game screen where all the things should happen
func GenerateGameScreen(sesh ssh.Session, sessions *map[ssh.Session]network.ActiveSession) (*tview.Grid, *tview.List) {
	var playerMessage string
	const LABEL string = "Message: "
	chatLog := tview.NewList()
	inputField := tview.NewForm().AddInputField(LABEL, "", 32, nil, func(text string) {
		playerMessage = text
	})
	gameScreen := tview.NewGrid().
		SetRows(3).
		SetColumns(30).
		AddItem(chatLog, 1, 0, 3, 3, 0, 0, false).
		SetBorders(true).
		AddItem(inputField, 0, 0, 1, 3, 0, 0, true)
	inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetDoneFunc(func(key tcell.Key) {
		println(key)
		if key == tcell.KeyEnter {
			if playerMessage == "" {
				return
			}
			switch input.ParseInput(playerMessage) {
			case input.Message:
				network.BroadcastMessage(sessions, playerMessage, sesh, func(message string, sender credentials.Account, currentSession network.ActiveSession) {
					AppendText(currentSession.Chat, sender, message, currentSession.UI)
				})
			}
			playerMessage = ""
			inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetText("")
			(*sessions)[sesh].UI.SetFocus(inputField.GetFormItemByLabel(LABEL))
		}
	})
	return gameScreen, chatLog
}
func GenerateWelcomeScreen(app *tview.Application, sesh ssh.Session, currentlyActiveSessions map[ssh.Session]network.ActiveSession, db *bolt.DB) *tview.Modal {
	welcome := tview.NewModal().
		SetText("Welcome to Akevitt. Would you like to register an account?").
		AddButtons([]string{"Yes", "Login"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Login" {
				app.SetRoot(generateLoginScreen(sesh, &currentlyActiveSessions, db), true)
			} else if buttonLabel == "Yes" {
				app.SetRoot(generateRegistrationScreen(sesh, &currentlyActiveSessions, db))
			}
		})
	return welcome
}
