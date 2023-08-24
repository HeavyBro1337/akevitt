package main

import (
	"akevitt"
	"bytes"
	"fmt"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

func registerScreen(engine *akevitt.Akevitt, session *TemplateSession) tview.Primitive {
	var username string
	var password string
	var repeatPassword string

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
			if password != repeatPassword {
				errorBox("Passwords don't match!", session, session.previousUI)
				return
			}
			err := engine.Register(username, password, session)
			if err != nil {
				errorBox(err.Error(), session, session.previousUI)
				return
			}
			session.SetRoot(characterCreationWizard(engine, session))
		}).
		AddButton("Back", func() {
			session.app.SetRoot(rootScreen(engine, session), true)
		})
	registerScreen.SetBorder(true).SetTitle(" Register ")
	return registerScreen
}

func characterCreationWizard(engine *akevitt.Akevitt, session *TemplateSession) tview.Primitive {
	var name string
	var description string
	characterCreator := tview.NewForm().AddInputField("Character Name: ", "", 32, nil, func(text string) {
		name = text
	}).AddTextArea("Character Description: ", "", 64, 64, 0, func(text string) {
		description = text
	})
	characterCreator.AddButton("Done", func() {
		if strings.TrimSpace(name) == "" {
			errorBox("character name must not be empty!", session, session.previousUI)
			return
		}
		characterParams := CharacterParams{}
		characterParams.name = name
		characterParams.description = description
		emptyChar := &Character{}

		_, err := akevitt.CreateObject(engine, session, emptyChar, characterParams)
		if err != nil {
			errorBox(err.Error(), session, session.previousUI)
			return
		}
		session.SetRoot(gameScreen(engine, session))
	})
	return characterCreator
}

func loginScreen(engine *akevitt.Akevitt, session *TemplateSession) tview.Primitive {
	var username string
	var password string
	loginScreen := tview.NewForm().
		AddInputField("Username: ", "", 32, nil, func(text string) {
			username = text
		}).
		AddPasswordField("Password: ", "", 32, '*', func(text string) {
			password = text
		})
	loginScreen.
		AddButton("Login", func() {
			err := engine.Login(username, password, session)
			if err != nil {
				errorBox(err.Error(), session, session.previousUI)
				return
			}
			character, err := akevitt.FindObject[*Character](engine, session, CharacterKey)

			if err != nil {
				errorBox(err.Error(), session, session.previousUI)
				session.SetRoot(characterCreationWizard(engine, session))
				return
			}
			session.character = character
			room, err := engine.GetRoom(session.character.CurrentRoomKey)

			if err != nil {
				errorBox(err.Error(), session, session.previousUI)
				return
			}
			session.character.account = session.account
			session.character.currentRoom = room
			room.AddObjects(session.character)
			session.SetRoot(gameScreen(engine, session))
		}).
		AddButton("Back", func() {
			session.app.SetRoot(rootScreen(engine, session), true)
		})
	return loginScreen
}

// Root screen is a screen which gets displayed when you connect via SSH.
// The root screen may lead to the authentication process and then gameplay screen.
func rootScreen(engine *akevitt.Akevitt, session akevitt.ActiveSession) tview.Primitive {
	sess, ok := session.(*TemplateSession)

	if !ok {
		panic("could not cast to custom session")
	}

	b, err := os.ReadFile("./data/logo.png")
	if err != nil {
		panic("Cannot find the image!")
	}
	pngLogo, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err.Error())
	}
	image := tview.NewImage().SetImage(pngLogo)
	wizard := tview.NewModal().
		SetText("Welcome to the template! Would you register your account?").
		AddButtons([]string{"Register", "Login"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Login" {
				sess.SetRoot(loginScreen(engine, sess))
			} else if buttonLabel == "Register" {
				sess.SetRoot(registerScreen(engine, sess))
			}
		})
	welcome := tview.NewGrid().
		SetBorders(false).
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		AddItem(image, 0, 0, 3, 27, 0, 0, false).
		AddItem(wizard, 2, 2, 3, 3, 0, 0, true)

	sess.app.SetFocus(wizard)
	return welcome
}

// Gameplay screen
func gameScreen(engine *akevitt.Akevitt, session *TemplateSession) tview.Primitive {
	playerMessage := ""

	// Preparing session by initializing UI primitives, channels and collections.
	chatlog := logview.NewLogView()
	chatlog.SetLevelHighlighting(true)
	session.subscribedChannels = []string{"ooc"}
	session.proceed = make(chan struct{})
	session.chat = chatlog

	inputField := tview.NewInputField().SetAutocompleteFunc(func(currentText string) (entries []string) {
		if len(currentText) == 0 {
			return
		}
		for _, word := range engine.GetCommands() {
			if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
				entries = append(entries, word)
			}
		}

		f, ok := autocompletion[strings.Split(currentText, " ")[0]]

		if ok {
			for _, word := range f(currentText, engine, session) {
				if strings.HasPrefix(strings.ToLower(word), strings.ToLower(currentText)) {
					entries = append(entries, word)
				}
			}
		}
		return entries
	}).SetChangedFunc(func(text string) {
		playerMessage = text
	})
	session.input = inputField
	// Creating some useful UI elements such as character's status (health, money, etc.)
	// and visible objects in a room.
	status := stats(engine, session)
	visibles := visibleObjects(engine, session)

	session.app.SetAfterDrawFunc(func(screen tcell.Screen) {
		lookupUpdate(engine, session, &visibles)
		fmt.Fprint(status.Clear(), updateStats(engine, session))
	})

	// The gamescreen to be returned
	gameScreen := tview.NewGrid().
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		AddItem(status, 0, 0, 1, 3, 0, 0, false).
		AddItem(visibles, 1, 0, 1, 1, 0, 0, false).
		AddItem(inputField, 2, 0, 1, 3, 0, 0, true).
		AddItem(chatlog, 1, 1, 1, 2, 0, 0, false).
		SetBorders(true)
	inputField.SetFinishedFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			playerMessage = strings.TrimSpace(playerMessage)
			if playerMessage == "" {
				inputField.SetText("")
				return
			}
			AppendText(session, "\t>"+playerMessage, session.chat)
			err := engine.ExecuteCommand(playerMessage, session)
			if err != nil {
				errorBox(err.Error(), session, session.previousUI)
				inputField.SetText("")
				return
			}
			playerMessage = ""
			inputField.SetText("")
			go func() {
				for {
					_, ok := <-session.proceed
					if !ok {
						break
					}
					session.app.SetFocus(inputField)
					lookupUpdate(engine, session, &visibles)
					fmt.Fprint(status.Clear(), updateStats(engine, session))
					session.character.Save(engine)
				}
			}()
		}
	})
	inputField.SetAutocompletedFunc(func(text string, index, source int) bool {
		if source != tview.AutocompletedNavigate {
			inputField.SetText(text)
		}
		return source == tview.AutocompletedEnter || source == tview.AutocompletedClick
	})

	return gameScreen
}

func dialogueBox(dial *akevitt.Dialogue, engine *akevitt.Akevitt, session *TemplateSession) error {
	labels := akevitt.MapSlice(dial.GetOptions(), func(v *akevitt.Dialogue) string {
		return v.GetTitle()
	})

	if len(labels) == 0 {
		session.app.SetRoot(*session.previousUI, true)
		session.proceed <- struct{}{}
		return nil
	}
	var err error = nil
	modal := tview.NewModal().AddButtons(labels).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonIndex < 0 {
			return
		}

		err = dial.Proceed(buttonIndex, session, engine)
		if len(dial.GetOptions()) == 0 {
			session.app.SetRoot(*session.previousUI, true)
		}
	}).SetText("Press escape to change focus")

	modal.SetBackgroundColor(tcell.ColorBlack).SetBorder(false)

	grid := func(p tview.Primitive) tview.Primitive {
		gr := tview.NewGrid().
			SetColumns(3).
			SetRows(3).
			AddItem(p, 0, 1, 3, 2, 0, 0, false).
			AddItem(modal, 1, 1, 1, 1, 0, 0, true)
		gr.SetBorder(true)
		gr.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			if event.Key() == tcell.KeyEscape {
				if session.app.GetFocus() == p {
					session.app.SetFocus(modal)
				} else {
					session.app.SetFocus(p)
				}
			}
			return event
		})
		return gr
	}

	session.app.SetRoot(grid(dial.GetContents()), true)
	return err
}

func errorBox(message string, session *TemplateSession, back *tview.Primitive) {
	result := tview.NewModal().SetText("Error!").SetText(message).SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		session.app.SetRoot(*back, true)
		if session.input != nil {
			session.app.SetFocus(session.input)
		}
	})

	result.SetBorder(true).SetBorderColor(tcell.ColorDarkRed)
	session.app.SetRoot(result, true)
}

func stats(engine *akevitt.Akevitt, session *TemplateSession) *tview.TextView {
	return tview.NewTextView().SetText(updateStats(engine, session))
}

func updateStats(engine *akevitt.Akevitt, session *TemplateSession) string {
	character := session.character
	return fmt.Sprintf("HEALTH: %d/%d, NAME: %s (%s) $%d", character.Health,
		character.MaxHealth,
		character.Name,
		character.currentRoom.GetName(),
		character.Money)
}

func visibleObjects(engine *akevitt.Akevitt, session *TemplateSession) *tview.List {
	l := tview.NewList()
	lookupUpdate(engine, session, &l)
	return l
}

func lookupUpdate(engine *akevitt.Akevitt, session *TemplateSession, l **tview.List) {
	objects := engine.Lookup(session.character.currentRoom)
	(*l).Clear()
	for _, v := range objects {
		if v == session.character {
			continue
		}

		(*l).AddItem(v.GetName(), v.GetDescription(), 0, nil)
	}
	(*l).AddItem("AVAILABLE ROOMS", "", 0, nil)
	exits := session.character.currentRoom.GetExits()

	for _, v := range exits {
		(*l).AddItem(v.GetRoom().GetName(), strconv.FormatUint(v.GetKey(), 10), 0, nil)
	}
	(*l).SetSelectedBackgroundColor(tcell.ColorBlack).SetSelectedTextColor(tcell.ColorWhite)
}

type ItemFunc = func(item Interactable)

func inventoryList[T Interactable](engine *akevitt.Akevitt, session *TemplateSession, f ItemFunc) *tview.List {
	l := tview.NewList()

	for _, v := range akevitt.FilterByType[T](session.character.Inventory) {
		vCopy := v

		l.AddItem(v.GetName(), v.GetDescription(), 0, func() {
			if f != nil {
				f(vCopy)
			}
		}).SetSelectedFunc(func(i int, s1, s2 string, r rune) {
			l.RemoveItem(i)
		})
	}

	return l
}

func AppendText(currentSession *TemplateSession, message string, chatlog *logview.LogView) {
	ev := logview.NewLogEvent("message", message)
	ev.Level = logview.LogLevelInfo
	chatlog.AppendEvent(ev)
	chatlog.SetFocusFunc(func() {
		chatlog.Blur()
	})
	chatlog.ScrollToBottom()
}
