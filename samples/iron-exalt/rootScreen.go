package ironexalt

import (
	"akevitt"
	"bytes"
	"image/png"
	"log"
	"os"

	"github.com/rivo/tview"
)

func rootScreen(engine *akevitt.Akevitt, session akevitt.ActiveSession) tview.Primitive {
	sess, ok := session.(*ActiveSession)

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
		SetText("Welcome to Iron Exalt! Would you register your account?").
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
