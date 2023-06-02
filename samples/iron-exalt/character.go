package main

import (
	"akevitt/akevitt"
	"errors"
	"fmt"
)

const currentCharacterKey string = "currentCharacter"

type Character struct {
	Name        string
	Health      int
	MaxHealth   int
	account     akevitt.Account
	currentRoom *Room
	Map         map[string]akevitt.Object
}

type CharacterParams struct {
	name string
}

func (character *Character) Create(engine *akevitt.Akevitt, session *akevitt.ActiveSession, params interface{}) error {
	characterParams, ok := params.(CharacterParams)

	if !ok {
		return errors.New("invalid params given")
	}

	character.Name = characterParams.name
	character.Health = 10
	character.MaxHealth = 10
	character.account = *session.Account
	character.currentRoom = engine.GetSpawnRoom().(*Room)
	character.Map = make(map[string]akevitt.Object, 0)

	key, err := engine.GetNewKey(false)

	if err != nil {
		return err
	}

	return character.Save(key, engine)
}

func (character *Character) Save(key uint64, engine *akevitt.Akevitt) error {
	return engine.SaveObject(character, key)
}

func (character *Character) GetMap() map[string]akevitt.Object {
	return character.Map
}

func (character *Character) Description() string {
	return fmt.Sprintf("%s, HP: %d/%d", character.Name, character.Health, character.MaxHealth)
}

func (character *Character) GetAccount() akevitt.Account {
	return character.account
}

func (character *Character) GetRoom() akevitt.Room {
	return character.currentRoom
}
