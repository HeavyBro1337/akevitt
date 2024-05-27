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
		UseSpawnRoom(&room).
		UseRootUI(func(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
			return tview.NewTextView().SetText("Hello, World!")
		}).
		UseBind(":2222").
		Finish()

	log.Fatal(app.Run())
}
