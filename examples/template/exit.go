package main

import (
	"akevitt"
	"errors"
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
	sess, ok := session.(*TemplateSession)
	if !ok {
		return errors.New("invalid session type")
	}
	character := sess.character
	character.currentRoom.RemoveObject(character)
	room := exit.room
	character.currentRoom = room
	room.AddObjects(character)

	character.CurrentRoomKey = exit.Key
	return character.Save(engine)
}

func (exit *Exit) GetRoom() akevitt.Room {
	return exit.room
}

func (exit *Exit) Save(engine *akevitt.Akevitt) error {
	return akevitt.SaveObject[*Exit](engine, exit, "Exits", exit.Key)
}

func (exit *Exit) SetRoom(room akevitt.Room) {
	exit.Key = room.GetKey()
	exit.room = room
}
