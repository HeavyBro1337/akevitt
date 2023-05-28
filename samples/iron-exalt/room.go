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
	Exits           []akevitt.Exit
	ObjectKeys      []uint64
	Key             uint64
}

type RoomParams struct {
	name  string
	exits []akevitt.Exit
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

func (room *Room) Save(key uint64, engine *akevitt.Akevitt) error {
	return engine.SaveWorldObject(room, key)
}

func (room *Room) GetKeys(engine *akevitt.Akevitt) ([]uint64, error) {
	return room.ObjectKeys, nil
}

func (room *Room) AddObject(engine *akevitt.Akevitt, key uint64) error {
	room.ObjectKeys = append(room.ObjectKeys, key)

	return room.Save(room.Key, engine)
}

func (room *Room) RemoveObject(engine *akevitt.Akevitt, key uint64) error {
	room.ObjectKeys = remove(room.ObjectKeys, key)

	return room.Save(room.Key, engine)
}

func remove[T comparable](l []T, item T) []T {
	out := make([]T, 0)
	for _, element := range l {
		if element != item {
			out = append(out, element)
		}
	}
	return out
}
