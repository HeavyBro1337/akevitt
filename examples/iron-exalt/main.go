// This is an example game powered by Akevitt named Iron Exalt
// This game implements gameplay elements such as dialogues, talking with NPCs, selling ores, etc.

package main

import (
	"akevitt"
	"encoding/gob"
	"log"
)

const (
	CharacterKey uint64 = iota + 1
	NpcKey
)

func main() {
	initAutocompletion()

	gob.Register(&BaseItem{}) // Registering custom structs for database
	gob.Register(&Exit{})
	gob.Register(&Room{})
	gob.Register(&Ore{})

	spawn := generateRooms()

	engine := akevitt.NewEngine(). // Engine creation stage
					UseDBPath("data/iron-exalt.db"). // Using custom database path
					UseOnMessage(onMessage).         // Installing hooks
					UseOnSessionEnd(onSessionEnd).   // to have control over logic
					UseOnDialogue(onDialogue).       //
					UseRegisterCommand("say", say).  // Registering commands
					UseRegisterCommand("ooc", ooc).
					UseRegisterCommand("enter", enter).
					UseRegisterCommand("interact", interact).
					UseRegisterCommand("backpack", backpack).
					UseRegisterCommand("look", look).
					UseRegisterCommand("mine", mine).
					UseRegisterCommand("attack",
			func(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
				return engine.SubscribeToHeartBeat(15, func() {
					engine.Message("ooc", "Attacked", session.GetAccount().Username, session)
				})
			}).
		UseNewHeartbeat(15).
		UseSpawnRoom(spawn).  // Setting spawn root room
		UseRootUI(rootScreen) // Passing root screen

	log.Fatal(akevitt.Run[*IronExaltSession](engine))
}
