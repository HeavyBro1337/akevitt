/*
Program written by Maxwell Jensen, Ivan Korchmit (c) 2023
Licensed under European Union Public Licence 1.2.
For more information, view the man page or README.md
*/

package main

import (
  "golang.org/x/term"
  "github.com/gliderlabs/ssh"
  "github.com/boltdb/bolt"
  "io"
  "log"
)

func main() {
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
    io.WriteString(sesh, "Your address: " + sesh.RemoteAddr().String() + "\n")

    // Main loop that runs for as long as a client is connected
    for {
      io.WriteString(sesh, "Your input for database entry #1 (n/N to quit):\n")
      input, err := uTerm.ReadLine()
      // There might be errors, like EOF (Ctrl + C)
      if err != nil {
        io.WriteString(sesh, err.Error() + "\n")
        sesh.Close()
      }
      // See what the input does
      switch parseInput(input) {
      case true: sesh.Close()
      default: continue
      }
    }
  })

  log.Fatal(ssh.ListenAndServe(":2222", nil))
}

// Entry point for all client input
// If `exit` is `true`, client ends SSH session
func parseInput(inp string) (exit bool) {
  // Check that the string is empty, otherwise see if its q/Q
  if len(inp) == 0 {
    return false
  } else if len(inp) == 1 && inp[0] == 'q' || inp[0] == 'Q' {
    return true
  }
  return false
}
