package main

import (
	"akevitt/akevitt"
	"bytes"
	"errors"
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func AppendText(currentSession akevitt.ActiveSession, senderName string, message string, sh rune) error {
	if currentSession.Chat == nil {
		return errors.New("chat log element is nil")
	}
	currentSession.Chat.InsertItem(0, senderName, message, sh, nil)
	currentSession.Chat.SetWrapAround(true)

	return nil
}

func loginScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	var username string
	var password string
	loginScreen := tview.NewForm().
		AddInputField("Username: ", "", 32, nil, func(text string) {
			username = text
		}).
		AddPasswordField("Password: ", "", 32, '*', func(text string) {
			password = text
		})
	loginScreen.AddButton("Login", func() {
		err := engine.Login(username, password, session)

		if err != nil {
			ErrorBox(err.Error(), session.UI, session.UIPrimitive)
			return
		}

		character, key, err := akevitt.FindObject[*Character](engine, session)

		if err != nil {
			session.SetRoot(characterCreationWizard(engine, session))
			return
		}

		character.account = *session.Account
		err = character.OnLoad(engine)

		if err != nil {
			fmt.Printf("err loading character: %v\n", err)
			return
		}

		session.SetRoot(gameScreen(engine, session))
		session.RelatedGameObjects[currentCharacterKey] = akevitt.Pair[uint64, akevitt.GameObject]{First: key, Second: character}

	})
	return loginScreen
}

func ErrorBox(message string, app *tview.Application, back *tview.Primitive) {
	result := tview.NewModal().SetText("Error!").SetText(message).SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.SetRoot(*back, false)
	}).SetFocus(0)

	result.SetBorderColor(tcell.ColorDarkRed)
	result.SetBorder(true)
	app.SetRoot(result, true)
}

func gameScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	fmt.Printf("session: %v\n", session.Account)
	var playerMessage string
	const LABEL string = "Message: "
	session.Chat = tview.NewList()
	inputField := tview.NewForm().AddInputField(LABEL, "", 32, nil, func(text string) {
		playerMessage = text
	})

	gameScreen := tview.NewGrid().
		SetRows(3).
		SetColumns(30).
		AddItem(session.Chat, 1, 0, 3, 3, 0, 0, false).
		SetBorders(true).
		AddItem(inputField, 0, 0, 1, 3, 0, 0, true)
	inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			playerMessage = strings.TrimSpace(playerMessage)
			if playerMessage == "" {
				return
			}

			err := engine.ProcessCommand(playerMessage, session)

			if err != nil {
				ErrorBox(err.Error(), session.UI, session.UIPrimitive)
			}

			playerMessage = ""
			inputField.GetFormItemByLabel(LABEL).(*tview.InputField).SetText("")
			session.UI.SetFocus(inputField.GetFormItemByLabel(LABEL))
		}
	})
	return gameScreen
}

func registerScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
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

				ErrorBox("Passwords don't match!", session.UI, session.UIPrimitive)
				return
			}
			err := engine.Register(username, password, session)
			if err != nil {
				ErrorBox(err.Error(), session.UI, session.UIPrimitive)
				return
			}
			session.SetRoot(characterCreationWizard(engine, session))
		})
	registerScreen.SetBorder(true).SetTitle(" Register ")
	return registerScreen
}

func rootScreen(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	b, err := os.ReadFile("./data/logo.png")
	if err != nil {
		panic("Cannot find image!!!")
	}
	pngLogo, err := png.Decode(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err.Error())
	}
	image := tview.NewImage().SetImage(pngLogo)
	wizard := tview.NewModal().
		SetText(fmt.Sprintf("Welcome to %s! Would you register your account?", engine.GetGameName())).
		AddButtons([]string{"Register", "Login"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Login" {
				session.SetRoot(loginScreen(engine, session))
			} else if buttonLabel == "Register" {
				session.SetRoot(registerScreen(engine, session))
			}
		})
	welcome := tview.NewGrid().
		SetBorders(false).
		SetRows(3, 0, 3).
		SetColumns(30, 0, 30).
		AddItem(image, 0, 0, 3, 27, 0, 0, false).
		AddItem(wizard, 2, 2, 3, 3, 0, 0, false)
	return welcome
}

func characterCreationWizard(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	var name string
	characterCreator := tview.NewForm().AddInputField("Character Name: ", "", 32, nil, func(text string) {
		name = text
	})
	characterCreator.AddButton("Done", func() {
		if strings.TrimSpace(name) == "" {
			ErrorBox("character name must not be empty!", session.UI, session.UIPrimitive)
			return
		}
		characterParams := CharacterParams{}
		characterParams.name = name

		emptychar := &Character{}
		_, err := akevitt.CreateObject(engine, session, emptychar, characterParams)

		if err != nil {
			ErrorBox(err.Error(), session.UI, session.UIPrimitive)
			return
		}

		session.SetRoot(gameScreen(engine, session))
	})
	return characterCreator
}
