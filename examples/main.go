package main

import (
	"log"

	"github.com/IvanKorchmit/akevitt/engine"
	"github.com/IvanKorchmit/akevitt/plugins"
	"github.com/rivo/tview"
)

func main() {
	room := engine.NewRoom("Example room")

	app := engine.NewEngine().
		AddPlugin(plugins.DefaultPlugins()...).
		AddPlugin(plugins.NewAccountPlugin()).
		AddPlugin(plugins.NewBoltPlugin[*engine.Account]("database.db")).
		UseSpawnRoom(room).
		UseRootUI(Root).
		UseBind(":1999").
		Finish()

	log.Fatal(app.Run())
}

func Root(eng *engine.Akevitt, session *engine.ActiveSession) tview.Primitive {
	modal := tview.NewModal().AddButtons([]string{"Go!"})

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		session.Application.SetRoot(plugins.RegistrationScreen(eng, session, func(eng *engine.Akevitt, session *engine.ActiveSession) tview.Primitive {
			return tview.NewTextView().SetText("Thank you!!")
		}), true)
	})

	return modal
}