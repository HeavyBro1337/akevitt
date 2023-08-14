package main

import (
	"akevitt/akevitt"
	"fmt"
	"log"

	"github.com/rivo/tview"
)

const (
	CharacterKey uint64 = iota + 1
)

func main() {
	engine := akevitt.NewEngine().
		UseMouse().
		UseDBPath("data/iron-exalt.db").
		UseMessage(func(engine *akevitt.Akevitt, session akevitt.ActiveSession, channel, message, username string) error {
			sess, ok := session.(*ActiveSession)

			st := fmt.Sprintf("[black:red]%s (%s) [black:white]%s", username, channel, message)

			if ok && sess.subscribedChannels != nil {
				if find[string](sess.subscribedChannels, channel) {
					return AppendText(sess, st, sess.chat)
				}
			} else if !ok {
				fmt.Printf("could not cast to session")
			} else {
				fmt.Print("unknown error")
			}

			return nil
		}).
		RegisterCommand("say", say).
		RegisterCommand("ooc", ooc).
		UseRootUI(rootScreen)

	log.Fatal(akevitt.Run[*ActiveSession](engine))
}

func find[T comparable](collection []T, value T) bool {
	for _, b := range collection {
		fmt.Printf("b: %v\n", b)
		if b == value {
			return true
		}
	}
	return false
}

func AppendText(currentSession *ActiveSession, message string, chatlog *tview.List) error {
	chatlog.AddItem(message, "", 0, nil)
	chatlog.SetCurrentItem(-1)
	chatlog.SetWrapAround(true)
	chatlog.ShowSecondaryText(false)
	return nil
}
