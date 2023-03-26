package main

import (
	"akevitt/akevitt"
	"errors"
	"fmt"
)

const currentCharacterKey string = "currentCharacter"

type Character struct {
	name      string
	health    int
	maxHealth int
	account   akevitt.Account
}

type CharacterParams struct {
	name string
}

func (character *Character) Create(engine *akevitt.Akevitt, session *akevitt.ActiveSession, params interface{}) error {
	characterParams, ok := params.(CharacterParams)

	if !ok {
		return errors.New("invalid params given")
	}
	character.name = characterParams.name
	character.health = 10
	character.maxHealth = 10
	character.account = *session.Account

	return nil
}

func (character *Character) Save(key uint64, engine *akevitt.Akevitt) error {
	return nil
}

func (character *Character) Description() string {
	return fmt.Sprintf("%s, HP: %d/%d", character.name, character.health, character.maxHealth)
}

func (character *Character) GetAccount() akevitt.Account {
	return character.account
}
