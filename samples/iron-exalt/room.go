package main

import (
	"akevitt/akevitt"
	"errors"
)

type Exit struct {
	Name string
}

func (entrance *Exit) Description() string {
	return "Hello!"
}

func (entrance *Exit) Save(key uint64, engine *akevitt.Akevitt) error {
	return engine.SaveWorldObject(entrance, key)
}

type Room struct {
	Name            string
	DescriptionData string
	Exits           []uint64
	Key             uint64
}

type RoomParams struct {
	name  string
	exits []uint64
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

func (room *Room) GetExits() []uint64 {
	return room.Exits
}

func (room *Room) Save(key uint64, engine *akevitt.Akevitt) error {
	return engine.SaveWorldObject(room, key)
}
