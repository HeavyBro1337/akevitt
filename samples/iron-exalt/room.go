package main

import (
	"akevitt/akevitt"
	"errors"
)

type Room struct {
	Name            string
	DescriptionData string
	Exits           []akevitt.Exit
	Key             uint64
}

type RoomParams struct {
	name  string
	exits []akevitt.Exit
}

func (room *Room) GetKey() uint64 {
	return room.Key
}

func (room *Room) Create(engine *akevitt.Akevitt, session *akevitt.ActiveSession, params interface{}) error {
	roomParams, ok := params.(RoomParams)

	if !ok {
		return errors.New("invalid params given")
	}

	room.Name = roomParams.name
	room.Exits = roomParams.exits
	key, err := engine.GetNewKey(true)
	if err != nil {
		return err
	}

	room.Key = key

	return room.Save(key, engine)
}

func (room *Room) Description() string {
	return room.DescriptionData
}

func (room *Room) GetExits() []akevitt.Exit {
	return room.Exits
}

func (room *Room) SetExits(exits ...akevitt.Exit) {
	room.Exits = exits
}

func (room *Room) Save(key uint64, engine *akevitt.Akevitt) error {
	return engine.SaveWorldObject(room, key)
}

func (exit *Exit) Description() string {
	return "Hello!"
}

func (exit *Exit) Save(key uint64, engine *akevitt.Akevitt) error {
	return engine.SaveWorldObject(exit, key)
}

type Exit struct {
	Key  uint64
	room akevitt.Room
}

func (exit *Exit) Enter(engine *akevitt.Akevitt, session *akevitt.ActiveSession) error {
	return nil
}

func (exit *Exit) GetRoom() akevitt.Room {
	return exit.room
}

func (room *Room) OnLoad(engine *akevitt.Akevitt) error {
	println("Invoked OnLoad in room")
	for _, v := range room.Exits {
		otherRoom, err := akevitt.GetObject[akevitt.Room](engine, v.GetKey(), true)

		if err != nil {
			return err
		}

		v.SetRoom(otherRoom)
	}
	return nil
}

func (exit *Exit) OnLoad(engine *akevitt.Akevitt) error {
	return nil
}

func (exit *Exit) GetKey() uint64 {
	return exit.Key
}

func (exit *Exit) SetRoom(room akevitt.Room) {
	exit.Key = room.GetKey()
	exit.room = room
}

func BindRooms[T akevitt.Exit](room akevitt.Room, sampleExit T, otherRooms ...akevitt.Room) {
	var exits []akevitt.Exit = make([]akevitt.Exit, 0)

	for _, v := range otherRooms {
		exit := sampleExit
		exit.SetRoom(v)
		exits = append(exits, exit)
	}

	room.SetExits(exits...)
}
