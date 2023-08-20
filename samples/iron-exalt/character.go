package main

import (
	"akevitt/akevitt"
	"errors"
	"fmt"
)

type Character struct {
	Name           string
	Health         int
	MaxHealth      int
	account        *akevitt.Account
	currentRoom    akevitt.Room
	Inventory      []akevitt.Interactable
	CurrentRoomKey uint64
}

func (character *Character) Create(engine *akevitt.Akevitt, session akevitt.ActiveSession, params interface{}) error {
	fmt.Println("Creating character...")

	characterParams, ok := params.(CharacterParams)
	if !ok {
		return errors.New("invalid params given")
	}
	sess, ok := session.(*ActiveSession)

	if !ok {
		return errors.New("invalid session type")
	}

	character.Name = characterParams.name
	character.Health = 10
	character.MaxHealth = 10
	character.currentRoom = engine.GetSpawnRoom()
	character.Inventory = make([]akevitt.Interactable, 0)
	character.account = sess.account
	sess.character = character
	room := engine.GetSpawnRoom()
	room.ContainObjects(sess.character)
	character.currentRoom = room
	character.CurrentRoomKey = character.currentRoom.GetKey()

	return character.Save(engine)
}

func (character *Character) Save(engine *akevitt.Akevitt) error {
	return engine.SaveGameObject(character, CharacterKey, character.account)
}

func (character *Character) GetDescription() string {
	format := `
	Health %d/%d
	`
	return fmt.Sprintf(format, character.Health, character.MaxHealth)
}

func (character *Character) GetName() string {
	return character.Name
}

type CharacterParams struct {
	name string
}
