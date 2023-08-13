package main

import (
	"akevitt/akevitt"
	"log"
)

const (
	CharacterKey uint64 = iota + 1
)

func main() {
	log.Fatal(akevitt.NewEngine().
		UseMouse().
		UseDBPath("data/iron-exalt.db").
		UseOnMessage(func(engine *akevitt.Akevitt, session akevitt.ActiveSession, channel, message string) error {
			sess, ok := session.(*ActiveSession)

			if ok && sess.subscribedChannels != nil {
				if find[string](sess.subscribedChannels, channel) {
					return AppendText(sess, message)
				}
			}

			return nil
		}).
		UseRootUI(rootScreen).
		Run(&ActiveSession{}))
}

func find[T comparable](collection []T, value T) bool {
	for _, b := range collection {
		if b == value {
			return true
		}
	}
	return false
}

func AppendText(currentSession *ActiveSession, message string) error {
	currentSession.chat.AddItem(message, "", 0, nil)
	currentSession.chat.SetCurrentItem(-1)
	currentSession.chat.SetWrapAround(true)
	currentSession.chat.ShowSecondaryText(false)
	return nil
}
