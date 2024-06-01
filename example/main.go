package main

import (
	"log"

	"github.com/IvanKorchmit/akevitt"
	"github.com/IvanKorchmit/akevitt/plugins"
	"github.com/rivo/tview"
)

func main() {
	room := akevitt.Room{
		Name: "Example room",
	}

	app := akevitt.NewEngine().
		AddPlugin(plugins.DefaultPlugins()...).
		AddPlugin(plugins.NewAccountPlugin()).
		AddPlugin(plugins.NewBoltPlugin[*akevitt.Account]("database.db")).
		UseSpawnRoom(&room).
		UseRootUI(Root).
		UseBind(":1999").
		Finish()

	log.Fatal(app.Run())
}

func Root(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	modal := tview.NewModal().AddButtons([]string{"Go!"})

	modal.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		session.Application.SetRoot(plugins.RegistrationScreen(engine, session, func(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
			return tview.NewTextView().SetText("Thank you!!")
		}), true)
	})

	return modal
}
