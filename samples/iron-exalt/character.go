package main

import (
	"akevitt/akevitt"
	"errors"
	"fmt"
)

const currentCharacterKey string = "currentCharacter"

type Character struct {
	CharacterName  string
	Health         int
	MaxHealth      int
	account        akevitt.Account
	currentRoom    *Room
	Map            map[string]akevitt.Object
	CurrentRoomKey uint64
}

type CharacterParams struct {
	name string
}

func (character *Character) Create(engine *akevitt.Akevitt, session *akevitt.ActiveSession, params interface{}) error {
	characterParams, ok := params.(CharacterParams)

	if !ok {
		return errors.New("invalid params given")
	}
	character.CharacterName = characterParams.name
	character.Health = 10
	character.MaxHealth = 10
	character.currentRoom = engine.GetSpawnRoom().(*Room)
	character.Map = make(map[string]akevitt.Object, 0)
	character.account = *session.Account
	character.CurrentRoomKey = character.currentRoom.Key
	key, err := engine.GetNewKey(false)

	session.RelatedGameObjects[currentCharacterKey] = akevitt.Pair[uint64, akevitt.GameObject]{First: key, Second: character}

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

func (character *Character) OnLoad(engine *akevitt.Akevitt) error {
	println("Invoked on load")

	room, err := akevitt.GetObject[*Room](engine, character.CurrentRoomKey, true)

	if err != nil {
		return err
	}

	character.currentRoom = room

	return nil
}

func (character *Character) Description() string {
	return fmt.Sprintf("%s, HP: %d/%d", character.CharacterName, character.Health, character.MaxHealth)
}

func (character *Character) Name() string {
	return character.CharacterName
}

func (character *Character) OnRoomLookup() uint64 {
	return character.CurrentRoomKey
}
