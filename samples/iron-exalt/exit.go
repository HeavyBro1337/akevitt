package main

import (
	"akevitt/akevitt"
)

type Exit struct {
	Key  uint64
	room akevitt.Room
}

func (exit *Exit) GetName() string {
	if exit.room == nil {
		return "Nowhere"
	}
	return exit.room.GetName()
}

func (exit *Exit) GetKey() uint64 {
	return exit.Key
}

func (exit *Exit) Enter(engine *akevitt.Akevitt, session akevitt.ActiveSession) error {
	panic("unimplemented")
}
func (exit *Exit) GetRoom() akevitt.Room {
	return exit.room
}

func (exit *Exit) Save(engine *akevitt.Akevitt) error {
	return akevitt.SaveObject[*Exit](engine, exit, "Exits", exit.Key)
}

func (exit *Exit) SetRoom(room akevitt.Room) {
	exit.room = room
}
