package main

import (
	"akevitt/akevitt"
)

type Room struct {
	Name            string
	DescriptionData string
	Exits           []akevitt.Exit
	Key             uint64
}

func (room *Room) GetKey() uint64 {
	return room.Key
}

func (room *Room) GetName() string {
	return room.Name
}
