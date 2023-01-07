/*
Program written by Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view LICENCE or README
*/

package main

import (
	"errors"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
	"io"
	"log"
	"strings"
)

const COMMAND_PREFIX string = "/"

type InputType int

const (
	Ignore InputType = iota
	Command
	Message
)

type ActiveSession struct {
	account *Account
	ui      *tview.Application
	chat    *tview.List
}

func main() {
	//var sessions = make(map[ssh.Session]*Account)
	// Open the database file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("akevitt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var sessions = make(map[ssh.Session]ActiveSession)
	// Open the SSH session with any clients who connect
	ssh.Handle(func(sesh ssh.Session) {
		screen, err := NewSessionScreen(sesh)
		if err != nil {
			fmt.Fprintln(sesh.Stderr(), "unable to create screen:", err)
			return
		}

		purgeDeadSessions(&sessions)

		app := tview.NewApplication().SetScreen(screen).EnableMouse(true)
		sessions[sesh] = ActiveSession{chat: nil, account: nil, ui: app}
		welcome := tview.NewModal().
			SetText("Welcome to Akevitt. Would you like to register an account?").
			AddButtons([]string{"Yes", "Login"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Login" {
					app.SetRoot(generateLoginScreen(sesh, &sessions, db), true)
				} else if buttonLabel == "Yes" {
					app.SetRoot(generateRegistrationScreen(sesh, &sessions, db))
				}
			})

		app.SetRoot(welcome, false)
		if err := app.Run(); err != nil {
			fmt.Fprintln(sesh.Stderr(), err)
			return
		}

		sesh.Exit(0)
	})

	log.Fatal(ssh.ListenAndServe(":1487", nil))
}

// Broadcasts message
func broadcastMessage(sessions *map[ssh.Session]ActiveSession, message string, session ssh.Session) error {
	for key, element := range *sessions {
		// The user is not authenticated
		if element.account == nil {
			continue
		}
		appendText(element.chat, *(*sessions)[session].account, message, element.ui)
		// element.ui.Draw()
		if key != session {
			element.ui.ForceDraw()
		}
	}
	return nil
}

// Entry point for all client input
func parseInput(inp string) (status InputType) {

	// Check that the string is empty, otherwise see if its q/Q
	if len(inp) == 0 {
		return Ignore
	} else if strings.HasPrefix(inp, COMMAND_PREFIX) {
		// We entered command
		return Command
	}
	return Message
}

// Iterates through all currently dead sessions by trying to send null character.
// If it gets an error, then we found the dead session and we purge them from active ones.
func purgeDeadSessions(sessions *map[ssh.Session]ActiveSession) {
	for k := range *sessions {

		_, err := io.WriteString(k, "\000")
		if err != nil {
			delete(*sessions, k)
		}
	}
}

func errorBox[T tview.Primitive](message string, app *tview.Application, back *T) *tview.Modal {
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
func generateGameScreen(sesh ssh.Session, sessions *map[ssh.Session]ActiveSession) (*tview.Grid, *tview.List) {
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
			switch parseInput(playerMessage) {
			case Message:
				broadcastMessage(sessions, playerMessage, sesh)
			case Command:
				{
					sendPrivateMessageToClient((*sessions)[sesh], "This is test")
				}

			}
			playerMessage = ""
			inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetText("")
			(*sessions)[sesh].ui.SetFocus(inputField.GetFormItemByLabel(LABEL))
		}
	})
	return gameScreen, chatLog
}

func sendPrivateMessageToClient(active ActiveSession, message string) {
	appendNoSenderText(active.chat, message, active.ui)
}

// generate login screen
func generateLoginScreen(sesh ssh.Session, sessions *map[ssh.Session]ActiveSession, db *bolt.DB) *tview.Form {
	var username string
	var password string

	gameScreen, chatLog := generateGameScreen(sesh, sessions)

	loginScreen := tview.NewForm().AddInputField("Username: ", "", 32, nil, func(text string) {
		username = text
	}).
		AddPasswordField("Password: ", "", 32, '*', func(text string) {
			password = text
		})

	loginScreen.
		AddButton("Login", func() {
			purgeDeadSessions(sessions)
			ok, acc := Login(username, password, db)
			if ok {
				if !checkCurrentLogin(*acc, sessions) {
					if active, ok := (*sessions)[sesh]; ok {
						active.account = acc
						active.chat = chatLog
						(*sessions)[sesh] = active
					}
					(*sessions)[sesh].ui.SetRoot(gameScreen, true)

				} else {
					(*sessions)[sesh].ui.SetRoot(errorBox("Somebody already logged in!", (*sessions)[sesh].ui, &loginScreen), false)
				}
			} else {
				(*sessions)[sesh].ui.SetRoot(errorBox("Wrong password or username!", (*sessions)[sesh].ui, &loginScreen), false)
			}

		})
	loginScreen.SetBorder(true).SetTitle("Login")
	return loginScreen
}

// generate register screen
func generateRegistrationScreen(sesh ssh.Session, sessions *map[ssh.Session]ActiveSession, db *bolt.DB) (*tview.Form, bool) {
	var username string
	var password string
	var repeatPassword string

	gameScreen, chatLog := generateGameScreen(sesh, sessions)

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
			purgeDeadSessions(sessions)
			if password != repeatPassword {
				(*sessions)[sesh].ui.SetRoot(errorBox("Password don't match", (*sessions)[sesh].ui, &registerScreen),
					false)
				return
			}
			if !doesAccountExists(username, db) {
				acc := Account{Username: username, Password: password}
				createAccount(db, acc)
				if active, ok := (*sessions)[sesh]; ok {
					active.account = &acc
					active.chat = chatLog
					(*sessions)[sesh] = active
					(*sessions)[sesh].ui.SetRoot(gameScreen, true)
				}
			} else {
				(*sessions)[sesh].ui.SetRoot(errorBox("Account already exists", (*sessions)[sesh].ui, &registerScreen), false)
			}
		})
	registerScreen.SetBorder(true).SetTitle("Register")
	return registerScreen, true
}

// Appends message to the end-user
func appendText(text *tview.List, account Account, message string, ui *tview.Application) error {
	if text == nil {
		return errors.New("TextView is nil")
	}
	text.InsertItem(0, account.Username, message, 'M', nil)
	text.SetWrapAround(true)
	return nil
}

func appendNoSenderText(text *tview.List, message string, ui *tview.Application) error {
	if text == nil {
		return errors.New("TextView is nil")
	}
	text.
		AddItem("", message, 'S', nil)

	return nil
}
