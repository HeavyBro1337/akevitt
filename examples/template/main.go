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

	gob.Register(&BaseItem{})
	gob.Register(&Exit{})
	gob.Register(&Room{})
	gob.Register(&Ore{})

	spawn := generateRooms()

	engine := akevitt.NewEngine().
		UseDBPath("data/template.db").
		UseOnMessage(onMessage).
		UseOnSessionEnd(onSessionEnd).
		UseOnDialogue(onDialogue).
		UseRegisterCommand("say", say).
		UseRegisterCommand("ooc", ooc).
		UseRegisterCommand("enter", enter).
		UseRegisterCommand("interact", interact).
		UseRegisterCommand("backpack", backpack).
		UseRegisterCommand("look", look).
		UseRegisterCommand("mine", mine).
		UseNewHeartbeat(15).
		UseSpawnRoom(spawn).
		UseRootUI(rootScreen)

	log.Fatal(akevitt.Run[*TemplateSession](engine))
}
