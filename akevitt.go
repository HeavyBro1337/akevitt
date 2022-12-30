package main

import (
  "golang.org/x/term"
  "github.com/gliderlabs/ssh"
  "github.com/boltdb/bolt"
  "fmt"
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

  ssh.Handle(func(s ssh.Session) {
    uTerm := term.NewTerminal(s, "£ ")
    io.WriteString(s, "Your address: " + s.RemoteAddr().String() + "\n")
    var input string
    for {

      io.WriteString(s, "Your input for database entry #1 (n/N to quit):\n")
      input, err = uTerm.ReadLine()
      if err != nil {
        io.WriteString(s, err.Error() + "\n")
        s.Close()
      }
      if len(input) == 0 {
        continue
      } else if input[0] == 'n' || input[0] == 'N' {
        break
      }

      err = db.Update(func(tx *bolt.Tx) error {
        b, bErr := tx.CreateBucketIfNotExists([]byte("Main"))
        if bErr != nil {
          return fmt.Errorf("create bucket: %s", err)
        }
        pErr := b.Put(make([]byte, 1), []byte(input))
        return pErr
      })
      if err != nil {
        log.Fatal(err)
      }

      db.View(func(tx *bolt.Tx) error {
        b := tx.Bucket([]byte("Main"))
        v := b.Get(make([]byte, 1))
        io.WriteString(s, "Value of entry #1: " + string(v) + "\n")
        return nil
      })
    }
  })

  log.Fatal(ssh.ListenAndServe(":2222", nil))
}
