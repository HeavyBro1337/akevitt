package main

import (
	"akevitt/akevitt"
)

type Room struct {
	Name             string
	DescriptionData  string
	exits            []akevitt.Exit
	Key              uint64
	containedObjects []akevitt.GameObject
}

func (room *Room) OnCreate() {
	room.containedObjects = make([]akevitt.GameObject, 0)
}

func (room *Room) GetObjects() []akevitt.GameObject {
	return room.containedObjects
}

func (room *Room) Description() string {
	return room.DescriptionData
}

func (room *Room) GetKey() uint64 {
	return room.Key
}

func (room *Room) GetExits() []akevitt.Exit {
	return room.exits
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) Save(engine *akevitt.Akevitt) error {
	return akevitt.SaveObject[*Room](engine, room, "Rooms", room.Key)
}

func (room *Room) SetExits(exits ...akevitt.Exit) {
	room.exits = exits
}
