package main

import (
	"github.com/IvanKorchmit/akevitt"
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

func (room *Room) GetDescription() string {
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

func (room *Room) AddObjects(objects ...akevitt.GameObject) {
	room.containedObjects = append(room.containedObjects, objects...)
}

func (room *Room) RemoveObject(object akevitt.GameObject) {
	room.containedObjects = akevitt.RemoveItem(room.containedObjects, object)
}
