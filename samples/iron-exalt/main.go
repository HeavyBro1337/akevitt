package main

import (
	"akevitt/akevitt"
	"encoding/gob"
	"errors"
	"fmt"
	"log"

	"github.com/uaraven/logview"
)

const (
	CharacterKey uint64 = iota + 1
	NpcKey
)

func main() {
	gob.Register(&Exit{})
	gob.Register(&Room{})
	room := generateRooms()

	engine := akevitt.NewEngine().
		UseDBPath("data/iron-exalt.db").
		UseMessage(func(engine *akevitt.Akevitt, session akevitt.ActiveSession, channel, message, username string) error {
			if session == nil {
				return errors.New("session is nil. Probably the dead one")
			}

			sess, ok := session.(*ActiveSession)

			st := fmt.Sprintf("%s (%s): %s", username, channel, message)

			if ok && sess.subscribedChannels != nil {
				if akevitt.Find[string](sess.subscribedChannels, channel) {
					AppendText(sess, st, sess.chat)
				} else if sess.character.currentRoom.GetName() == channel {
					AppendText(sess, st, sess.chat)
				}
			} else if !ok {
				fmt.Printf("could not cast to session")
			} else {
				fmt.Print("unknown error")
			}

			return nil
		}).
		UseOnSessionEnd(func(deadSession akevitt.ActiveSession, liveSessions []akevitt.ActiveSession, engine *akevitt.Akevitt) {
			sess, ok := deadSession.(*ActiveSession)
			if !ok {
				fmt.Println("could not cast to session")
				return
			}
			if sess.account == nil {
				return
			}

			sess.character.currentRoom.RemoveObject(sess.character)
			for _, v := range liveSessions {
				lsess, ok := v.(*ActiveSession)

				if !ok || lsess.chat == nil {
					continue
				}

				AppendText(lsess, fmt.Sprintf("%s left the game", sess.account.Username), lsess.chat)
			}
		}).
		RegisterCommand("say", say).
		RegisterCommand("ooc", ooc).
		RegisterCommand("enter", enter).
		UseSpawnRoom(room).
		UseRootUI(rootScreen)

	log.Fatal(akevitt.Run[*ActiveSession](engine))
}

func generateRooms() *Room {

	room := &Room{
		Name:             "Spawn Room",
		DescriptionData:  "Just a spawn room.",
		exits:            []akevitt.Exit{},
		Key:              0,
		containedObjects: []akevitt.GameObject{},
	}
	room.ContainObjects(createNpc("Maxwell", "Jensen", 0))

	rooms := []akevitt.Room{
		room,
		&Room{Name: "Mine", DescriptionData: "Mine of the corporation.", Key: 1},
		&Room{Name: "Iron City", DescriptionData: "The lounge of the miners.", Key: 2},
	}
	akevitt.BindRooms[*Exit](room, rooms...)
	akevitt.BindRooms[*Exit](rooms[1], rooms...)
	akevitt.BindRooms[*Exit](rooms[2], rooms...)
	return room
}

func AppendText(currentSession *ActiveSession, message string, chatlog *logview.LogView) {
	ev := logview.NewLogEvent("message", message)
	ev.Level = logview.LogLevelInfo
	chatlog.AppendEvent(ev)
	chatlog.SetFocusFunc(func() {
		chatlog.Blur()
	})
}
