/*
Program written by Maxwell Jensen, Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view the man page or README.md
*/

package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/fatih/color"
	"github.com/gdamore/tcell/v2"
	"github.com/gliderlabs/ssh"
	"github.com/rivo/tview"
	"golang.org/x/term"
)

const COMMAND_PREFIX string = "/"

const LOGO string = `
   _____   __               .__  __    __   
  /  _  \ |  | __ _______  _|__|/  |__/  |_ 
 /  /_\  \|  |/ // __ \  \/ /  \   __\   __\
/    |    \    <\  ___/\   /|  ||  |  |  |  
\____|__  /__|_ \\___  >\_/ |__||__|  |__|  
        \/     \/    \/                    
`

type InputType int

const (
	Ignore InputType = iota
	Command
	Message
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
	var sessions = make(map[ssh.Session]*Account)
	// Open the SSH session with any clients who connect
	ssh.Handle(func(sesh ssh.Session) {
		sessions[sesh] = nil
		screen, err := NewSessionScreen(sesh)
		if err != nil {
			fmt.Fprintln(sesh.Stderr(), "unable to create screen:", err)
			return
		}


    purgeDeadSessions(&sessions)
		
    
    app := tview.NewApplication().SetScreen(screen).EnableMouse(true)
		var username string
		var password string
		var playerMessage string
		gameScreen := tview.NewForm().AddInputField("Say: ", "", 128, nil, func(text string) {
			playerMessage = text
		}).
			AddButton("Send", func() {
				broadcastMessage(&sessions, playerMessage, sesh)
			})
		loginScreen := tview.NewForm().AddInputField("Username: ", "", 32, nil, func(text string) {
			username = text
			println(username)
		}).
			AddPasswordField("Password: ", "", 32, '*', func(text string) {
				password = text
				println(password)
			})

      loginScreen.
      AddButton("Login", func() {
        purgeDeadSessions(&sessions)
				ok, acc := Login(username, password, db)
				if ok {
					if !checkCurrentLogin(*acc, &sessions) {
						sessions[sesh] = acc
						app.SetRoot(gameScreen, true)
					} else {
						app.SetRoot(errorBox("Somebody already logged in!", app, &loginScreen), false)
					}
				} else {
					app.SetRoot(errorBox("Wrong password or username!", app, &loginScreen), false)
				}

			})
    
		welcome := tview.NewModal().
			SetText("Welcome to Akevitt. Would you like to register an account?").
			AddButtons([]string{"Yes", "Login"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Login" {
					app.SetRoot(loginScreen, true)
				}
			})
      
		loginScreen.SetBorder(true).SetTitle("Login")
		app.SetRoot(welcome, false)
		if err := app.Run(); err != nil {
			fmt.Fprintln(sesh.Stderr(), err)
			return
		}

		sesh.Exit(0)
	})

	log.Fatal(ssh.ListenAndServe(":1487", nil))
}

// Performs login based on user input.
func loginCommand(sesh ssh.Session, uTerm *term.Terminal, db *bolt.DB, sessions map[ssh.Session]*Account) {
	name, err := uTerm.ReadLine()
	if err != nil {
		return
	}
	password, err := uTerm.ReadPassword("Enter password: ")
	if err != nil {
		return
	}
	exists, acc := Login(name, password, db)
	// We check if we logged in successfully
	if exists {
		// We also check if it was logged in already from the other user.
		if !checkCurrentLogin(*acc, &sessions) {
			sessions[sesh] = acc
			io.WriteString(sesh, fmt.Sprintf(color.GreenString("Login successful. Welcome back, %s\n"), acc.Username))
			return
		} else {
			io.WriteString(sesh, color.RedString("Error: this account is already logged in.\n"))
			return
		}
	} else {
		io.WriteString(sesh, color.RedString("Fail: wrong password or username.\n"))
		return
	}
}

// Broadcasts message
func broadcastMessage(sessions *map[ssh.Session]*Account, message string, session ssh.Session) error {
	for k, element := range *sessions {
		if k == session {
			// prevent broadcasting message to sender.
			if element == nil {
				io.WriteString(k, "Please "+color.GreenString("/login")+" or "+color.GreenString("/register")+"\n")
				continue
			}
			continue
		}
		// The user is not authenticated
		if element == nil {
			continue
		}
		if (*sessions)[session] == nil {
			continue
		}
		_, err := io.WriteString(k, fmt.Sprintf(color.CyanString("%s: %s\n"), (*sessions)[session].Username, message))
		if err != nil {
			delete(*sessions, k)
			continue
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
func removeSession(s *[]ssh.Session, i int) {
	(*s)[i] = (*s)[len(*s)-1]
	(*s) = (*s)[:len(*s)-1]
}

// Iterates through all currently dead sessions by trying to send null character.
// If it gets an error, then we found the dead session and we purge them from active ones.
func purgeDeadSessions(sessions *map[ssh.Session]*Account) {
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
