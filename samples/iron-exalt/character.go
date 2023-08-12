package main

import "akevitt/akevitt"

type Character struct {
	Name           string
	Health         int
	MaxHealth      int
	account        *akevitt.Account
	currentRoom    akevitt.Room
	Map            map[string]akevitt.Object
	CurrentRoomKey uint64
}

type CharacterParams struct {
	name string
}
