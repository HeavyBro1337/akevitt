/*
Program written by Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view LICENCE or README
*/

package main

import (
	"akevitt/core/network"
	"akevitt/core/ui"
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
)

func main() {
	//var sessions = make(map[ssh.Session]*Account)
	// Open the database file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("akevitt.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var currentlyActiveSessions = make(map[ssh.Session]network.ActiveSession)
	// Open the SSH session with any clients who connect
	ssh.Handle(func(sesh ssh.Session) {
		screen, err := ui.NewSessionScreen(sesh)
		if err != nil {
			fmt.Fprintln(sesh.Stderr(), "unable to create screen:", err)
			return
		}

		network.PurgeDeadSessions(&currentlyActiveSessions)

		app := tview.NewApplication().SetScreen(screen).EnableMouse(true)
		currentlyActiveSessions[sesh] = network.ActiveSession{Chat: nil, Account: nil, UI: app}
		welcome := ui.GenerateWelcomeScreen(app, sesh, currentlyActiveSessions, db)

		app.SetRoot(welcome, false)
		if err := app.Run(); err != nil {
			fmt.Fprintln(sesh.Stderr(), err)
			return
		}

		sesh.Exit(0)
	})

	log.Fatal(ssh.ListenAndServe(":2222", nil))
}
