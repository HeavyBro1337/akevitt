package main

import (
	"akevitt/akevitt"
	"errors"
)

type Entrance struct {
	Name string
}

func (entrance *Entrance) Description() string {
	return "Hello!"
}

func (entrance *Entrance) Save(key uint64, engine *akevitt.Akevitt) error {
	return engine.SaveWorldObject(entrance, key)
}

type Room struct {
	Name            string
	DescriptionData string
	Entrances       []akevitt.Entrance
}

type RoomParams struct {
	name     string
	entrance []akevitt.Entrance
}

func (room *Room) Create(engine *akevitt.Akevitt, session *akevitt.ActiveSession, params interface{}) error {
	roomParams, ok := params.(RoomParams)

	if !ok {
		return errors.New("invalid params given")
	}

	room.Name = roomParams.name
	room.Entrances = roomParams.entrance
	key, err := engine.GetNewKey(true)

	if err != nil {
		return err
	}

	return room.Save(key, engine)
}

func (room Room) Description() string {
	return room.DescriptionData
}

func (room Room) GetEntrances() []akevitt.Entrance {
	return room.Entrances
}

func (room Room) Save(key uint64, engine *akevitt.Akevitt) error {
	return engine.SaveWorldObject(room, key)
}
