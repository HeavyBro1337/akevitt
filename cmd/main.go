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
	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)
const COMMAND_PREFIX string = "/"

type InputType int
const (
  Ignore InputType = iota
  Command
  Message
)
func main() {
  var sessions = make(map[ssh.Session]*Account)
  // Open the database file in your current directory.
  // It will be created if it doesn't exist.
  db, err := bolt.Open("akevitt.db", 0600, nil)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  // Open the SSH session with any clients who connect
  ssh.Handle(func(sesh ssh.Session) {
    // Open VT100 terminal for client
    uTerm := term.NewTerminal(sesh, "£ ")
    sessions[sesh] = nil
    io.WriteString(sesh, "Your address: " + sesh.RemoteAddr().String() + "\n")
    // Main loop that runs for as long as a client is connected
    for {
      input, err := uTerm.ReadLine()
      // There might be errors, like EOF (Ctrl + C)
      if err != nil {
        io.WriteString(sesh, err.Error() + "\n")
        err = sesh.Close()
        if err != nil {
          break
        }
      }
      // See what the input does
      switch parseInput(input) {
      case Command: {
        input = strings.ToLower(input) // Make commands case-insensitive. Requires discussion about that!
        input = input[1:] // Removing command prefix
        fmt.Println(input)
        if input == "q" || input == "exit" {
          sesh.Close()
        } else if input == "register" { // Todo: Wrap command handling into a separate function.
          io.WriteString(sesh,"\nEnter username: ")
          name, err := uTerm.ReadLine()
          if err != nil {
            break
          }
          password, err := uTerm.ReadPassword("Enter password: ")
          if err != nil {
            break
          }
          acc := Account{Username: name, Password: password}
          id, err := createAccount(db,acc)
          // Account does already exist
          if id == 0 {
            io.WriteString(sesh, color.RedString("Error: this user does already exist.\n"))
            continue
          }
          if err != nil {
            break 
          }
          acc, err = getAccount(id, db)
          sessions[sesh] = &acc
          if err != nil {
            fmt.Printf("Error: %s", err)
            break
          }
          io.WriteString(sesh, color.GreenString("Successfully created account\n"))
        } else if input == "login" {
          name, err := uTerm.ReadLine()
          if err != nil {
            break
          }
          password, err := uTerm.ReadPassword("Enter password: ")
          if err != nil {
            break
          }
          exists, acc := Login(name, password, db)
          // We check if we logged in successfully
          if exists {
            // We also check if it was logged in already from the other user.
            if !checkCurrentLogin(*acc,&sessions) {
            sessions[sesh] = acc
            io.WriteString(sesh, fmt.Sprintf(color.GreenString("Login successful. Welcome back, %s\n"), acc.Username))
            } else {
              io.WriteString(sesh, color.RedString("Error: this account is already logged in.\n"));
            }
          } else {
            io.WriteString(sesh, color.RedString("Fail: wrong password or username.\n"))
          }
        } else {
          io.WriteString(sesh, color.RedString("Error! Unknown command.\n"))
          }
      }
    case Message: broadcastMessage(&sessions, input, sesh)
      }
    }
  })

  log.Fatal(ssh.ListenAndServe(":1487", nil))
}

// Broadcasts message
func broadcastMessage(sessions *map[ssh.Session]*Account, message string, session ssh.Session ) error {
  for k, element := range *sessions {
    if k == session {
      // prevent broadcasting message to sender.
      if element == nil {
        io.WriteString(k, "Please /login or /register\n")
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
    _, err := io.WriteString(k,fmt.Sprintf(color.CyanString("%s: %s\n"), (*sessions)[session].Username, message))
    if err != nil {
      delete(*sessions, k)
      continue
    }
  }
  return nil
}

// Entry point for all client input
// If `exit` is `true`, client ends SSH session
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