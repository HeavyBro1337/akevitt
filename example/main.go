package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/IvanKorchmit/akevitt"
	"github.com/IvanKorchmit/akevitt/plugins"
	"github.com/IvanKorchmit/akevitt/ui"
	"github.com/rivo/tview"
)

func main() {
	room := akevitt.Room{
		Name: "Example room",
	}

	app := akevitt.NewEngine().
		AddPlugin(plugins.DefaultPlugins()...).
		AddPlugin(plugins.NewBoltPlugin("database.db")).
		UseSpawnRoom(&room).
		UseRootUI(Root).
		UseMouse().
		UseBind(":1999").
		Finish()

	log.Fatal(app.Run())
}

func Root(engine *akevitt.Akevitt, session *akevitt.ActiveSession) tview.Primitive {
	// Loading HTML. Can be simplified or be a separate function, whatever.
	f, _ := os.ReadFile("./grid.html")
	r := bytes.NewReader(f)
	templ := template.New("html")
	builder := &strings.Builder{}
	io.Copy(builder, r)
	templ.Parse(builder.String())
	builder.Reset()
	// Executing templates
	templ.Execute(builder, map[string]any{
		"randomNumber": rand.Int(),
	})
	dom, _ := ui.ParseHTML(strings.NewReader(builder.String()), session.Application)
	// Registering callbacks
	dom.OnCallback("exit", func() {
		session.Kill(engine)
	})

	return dom.Primitive()
}
