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

func (exit *Exit) Enter(engine *akevitt.Akevitt, session *ActiveSession) error {
	panic("unimplemented")
}
func (exit *Exit) GetRoom() akevitt.Room {
	return exit.room
}
