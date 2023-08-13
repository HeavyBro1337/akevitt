package main

import (
	"akevitt/akevitt"
	"errors"
)

type Character struct {
	Name           string
	Health         int
	MaxHealth      int
	account        *akevitt.Account
	currentRoom    akevitt.Room
	CurrentRoomKey uint64
}

func (character *Character) Create(engine *akevitt.Akevitt, session akevitt.ActiveSession, params interface{}) error {
	characterParams, ok := params.(CharacterParams)
	if !ok {
		return errors.New("invalid params given")
	}
	character.Name = characterParams.name
	character.Health = 10
	character.MaxHealth = 10
	character.currentRoom = engine.GetSpawnRoom()
	character.account = session.GetAccount()
	// character.CurrentRoomKey = character.currentRoom.GetKey()

	return character.Save(engine)
}

func (character *Character) Save(engine *akevitt.Akevitt) error {
	return engine.SaveGameObject(character, CharacterKey, character.account)
}

type CharacterParams struct {
	name string
}
