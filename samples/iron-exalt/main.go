package main

import (
	"akevitt/akevitt"
	"encoding/gob"
	"errors"
	"fmt"
	"log"

	"github.com/rivo/tview"
	"github.com/uaraven/logview"
)

const (
	CharacterKey uint64 = iota + 1
	NpcKey
)

func main() {
	autocompletion["interact"] = func(entry string, engine *akevitt.Akevitt, session *ActiveSession) []string {
		npcs := akevitt.LookupOfType[*NPC](session.character.currentRoom)

		return akevitt.MapSlice[*NPC, string](npcs, func(v *NPC) string {
			return "interact " + v.Name
		})
	}

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
		UseDialogue(func(engine *akevitt.Akevitt, session akevitt.ActiveSession, dialogue *akevitt.Dialogue) error {
			sess, ok := session.(*ActiveSession)

			if !ok {
				return errors.New("could not cast to session")
			}

			err := DialogueBox(dialogue, engine, sess)
			return err
		}).
		RegisterCommand("say", say).
		RegisterCommand("ooc", ooc).
		RegisterCommand("enter", enter).
		RegisterCommand("interact", interact).
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
	room.ContainObjects(createNpc("Maxwell Jensen", "The tutor", 0).UseInteract(func(engine *akevitt.Akevitt, session *ActiveSession) error {
		lore := tview.NewTextView().SetText(`
		Welcome, fellow miner! In our corporation you will mine minerals and ores for us!
		Go to the mine and mine with your pickaxe. 
		Come back and ask Ivan to deposite ores and receive the salary.
		Good luck!
		`)
		cya := tview.NewTextView().SetText(`See you later!`)
		d := akevitt.NewDialogue("Welcome")
		d.SetContent(lore)
		d.
			AddOption("Ok", cya).
			End()
		return engine.Dialogue(d, session)
	}))
	room.ContainObjects(createNpc("Ivan Korchmit", "Depositor", 1))
	room.ContainObjects(createNpc("John Doe", "Merchant", 2))

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
