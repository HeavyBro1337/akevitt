package basic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/IvanKorchmit/akevitt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

func AppendText(currentSession *Session, message string, chatlog *logview.LogView) {
	ev := logview.NewLogEvent("message", message)
	ev.Level = logview.LogLevelInfo
	chatlog.AppendEvent(ev)
	chatlog.SetFocusFunc(func() {
		chatlog.Blur()
	})
	chatlog.ScrollToBottom()
}

func DialogueBox(dial *akevitt.Dialogue, engine *akevitt.Akevitt, session *Session) error {
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

func ErrorBox(message string, session *Session, back *tview.Primitive) {
	result := tview.NewModal().SetText("Error!").SetText(message).SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		session.app.SetRoot(*back, true)
		if session.Input != nil {
			session.app.SetFocus(session.Input)
		}
	})

	result.SetBorder(true).SetBorderColor(tcell.ColorDarkRed)
	session.app.SetRoot(result, true)
}

func registerScreen(engine *akevitt.Akevitt, session *Session, gameName string, gameScreen func(engine *akevitt.Akevitt, session *Session) tview.Primitive) tview.Primitive {
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
				ErrorBox("Passwords don't match!", session, session.GetCurrentUI())
				return
			}
			err := engine.Register(username, password, session)
			if err != nil {
				ErrorBox(err.Error(), session, session.GetCurrentUI())
				return
			}
			session.SetRoot(characterCreationWizard(engine, session, gameScreen))
		}).
		AddButton("Back", func() {
			session.app.SetRoot(RootScreen(engine, session, gameName, gameScreen), true)
		})
	registerScreen.SetBorder(true).SetTitle(" Register ")
	return registerScreen
}

func characterCreationWizard(engine *akevitt.Akevitt, session *Session, gameScreen func(engine *akevitt.Akevitt, session *Session) tview.Primitive) tview.Primitive {
	var name string
	var description string
	characterCreator := tview.NewForm().AddInputField("Character Name: ", "", 32, nil, func(text string) {
		name = text
	}).AddTextArea("Character Description: ", "", 64, 64, 0, func(text string) {
		description = text
	})
	characterCreator.AddButton("Done", func() {
		if strings.TrimSpace(name) == "" {
			ErrorBox("character name must not be empty!", session, session.previousUI)
			return
		}
		characterParams := CharacterParams{}
		characterParams.name = name
		characterParams.description = description
		emptyChar := &Character{}

		_, err := akevitt.CreateObject(engine, session, emptyChar, characterParams)
		if err != nil {
			ErrorBox(err.Error(), session, session.previousUI)
			return
		}
		session.SetRoot(gameScreen(engine, session))
	})
	return characterCreator
}

func GameScreen(engine *akevitt.Akevitt, session *Session) tview.Primitive {
	playerMessage := ""

	// Preparing session by initializing UI primitives, channels and collections.
	chatlog := logview.NewLogView()
	chatlog.SetLevelHighlighting(true)
	session.subscribedChannels = []string{"ooc"}
	session.proceed = make(chan struct{})
	session.Chat = chatlog

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
	session.Input = inputField
	// Creating some useful UI elements such as character's status (health, money, etc.)
	// and visible objects in a room.
	status := stats(engine, session)
	visibles := VisibleObject(engine, session)

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
			AppendText(session, "\t>"+playerMessage, session.Chat)
			err := engine.ExecuteCommand(playerMessage, session)
			if err != nil {
				ErrorBox(err.Error(), session, session.previousUI)
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
					session.Character.Save(engine)
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

func loginScreen(engine *akevitt.Akevitt, session *Session, gameName string, gameScreen func(engine *akevitt.Akevitt, session *Session) tview.Primitive) tview.Primitive {
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
				ErrorBox(err.Error(), session, session.previousUI)
				return
			}
			character, err := akevitt.FindObject[*Character](engine, session, CharacterKey)

			if err != nil {
				session.SetRoot(characterCreationWizard(engine, session, gameScreen))
				ErrorBox(err.Error(), session, session.previousUI)
				return
			}
			session.Character = character
			room, err := engine.GetRoom(session.Character.CurrentRoomKey)

			if err != nil {
				ErrorBox(err.Error(), session, session.previousUI)
				return
			}
			session.Character.currentRoom = room
			room.AddObjects(session.Character)
			session.SetRoot(GameScreen(engine, session))
		}).
		AddButton("Back", func() {
			session.app.SetRoot(RootScreen(engine, session, gameName, gameScreen), true)
		})
	return loginScreen
}

func RootScreen(engine *akevitt.Akevitt, session akevitt.ActiveSession, gameName string,
	gameScreen func(engine *akevitt.Akevitt, session *Session) tview.Primitive) tview.Primitive {
	sess := CastSession[*Session](session)

	wizard := tview.NewModal().
		SetText(fmt.Sprintf("Welcome to the %s! Would you register your account?", gameName)).
		AddButtons([]string{"Register", "Login"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Login" {
				sess.SetRoot(loginScreen(engine, sess, gameName, gameScreen))
			} else if buttonLabel == "Register" {
				sess.SetRoot(registerScreen(engine, sess, gameName, gameScreen))
			}
		})
	welcome := tview.NewGrid().
		SetBorders(false).
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		AddItem(wizard, 2, 2, 3, 3, 0, 0, true)

	sess.app.SetFocus(wizard)
	return welcome
}

func stats(engine *akevitt.Akevitt, session *Session) *tview.TextView {
	return tview.NewTextView().SetText(updateStats(engine, session))
}

func updateStats(engine *akevitt.Akevitt, session *Session) string {
	character := session.Character
	return fmt.Sprintf("HEALTH: %d/%d, NAME: %s (%s) $%d", character.Health,
		character.MaxHealth,
		character.Name,
		character.currentRoom.GetName(),
		character.Money)
}

func lookupUpdate(engine *akevitt.Akevitt, session *Session, l **tview.List) {
	(*l).Clear()
	for _, v := range session.Character.currentRoom.GetObjects() {
		if v == session.Character {
			continue
		}

		(*l).AddItem(v.GetName(), v.GetDescription(), 0, nil)
	}
	(*l).AddItem("AVAILABLE ROOMS", "", 0, nil)
	exits := session.Character.currentRoom.GetExits()

	for _, v := range exits {
		(*l).AddItem(v.GetRoom().GetName(), strconv.FormatUint(v.GetKey(), 10), 0, nil)
	}
	(*l).SetSelectedBackgroundColor(tcell.ColorBlack).SetSelectedTextColor(tcell.ColorWhite)
}

func VisibleObject(engine *akevitt.Akevitt, session *Session) *tview.List {
	l := tview.NewList()
	lookupUpdate(engine, session, &l)
	return l
}
