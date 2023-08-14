package main

import (
	"akevitt/akevitt"
	"encoding/gob"
	"fmt"
	"log"

	"github.com/uaraven/logview"
)

const (
	CharacterKey uint64 = iota + 1
)

func main() {
	gob.Register(&Exit{})
	gob.Register(&Room{})
	room := &Room{Name: "Spawn Room", DescriptionData: "Just a spawn room.", Key: 0}
	rooms := []akevitt.Room{
		room,
		&Room{Name: "Mine", DescriptionData: "Mine of the corporation.", Key: 1},
		&Room{Name: "Iron City", DescriptionData: "The lounge of the miners.", Key: 2},
	}
	akevitt.BindRooms[*Exit](room, rooms...)
	akevitt.BindRooms[*Exit](rooms[1], rooms...)
	akevitt.BindRooms[*Exit](rooms[2], rooms...)

	engine := akevitt.NewEngine().
		UseDBPath("data/iron-exalt.db").
		UseMessage(func(engine *akevitt.Akevitt, session akevitt.ActiveSession, channel, message, username string) error {
			sess, ok := session.(*ActiveSession)

			st := fmt.Sprintf("%s (%s): %s", username, channel, message)

			if ok && sess.subscribedChannels != nil {
				if akevitt.Find[string](sess.subscribedChannels, channel) {
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
		UseSpawnRoom(room).
		UseRootUI(rootScreen)

	log.Fatal(akevitt.Run[*ActiveSession](engine))
}

func AppendText(currentSession *ActiveSession, message string, chatlog *logview.LogView) error {
	ev := logview.NewLogEvent("message", message)
	ev.Level = logview.LogLevelInfo
	chatlog.AppendEvent(ev)
	chatlog.SetFocusFunc(func() {
		chatlog.Blur()
	})
	return nil
}
