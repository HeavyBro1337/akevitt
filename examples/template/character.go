package main

import (
	"akevitt"
	"errors"
	"fmt"
)

// The character struct represents an in-game character that the user will play as.
type Character struct {
	Name           string
	Description    string
	Health         int
	Money          int
	MaxHealth      int
	account        *akevitt.Account
	currentRoom    akevitt.Room
	Inventory      []Interactable
	CurrentRoomKey uint64
}

func (character *Character) Create(engine *akevitt.Akevitt, session akevitt.ActiveSession, params interface{}) error {
	fmt.Println("Creating a character...")

	characterParams, ok := params.(CharacterParams)
	if !ok {
		return errors.New("invalid params given")
	}
	sess, ok := session.(*TemplateSession)

	if !ok {
		return errors.New("invalid session type")
	}

	character.Name = characterParams.name
	character.Description = characterParams.description
	character.Health = 10
	character.MaxHealth = 10
	character.Money = 100
	character.currentRoom = engine.GetSpawnRoom()
	character.Inventory = make([]Interactable, 0)
	character.account = sess.account
	sess.character = character
	room := engine.GetSpawnRoom()
	room.AddObjects(sess.character)
	character.currentRoom = room
	character.CurrentRoomKey = character.currentRoom.GetKey()

	pick := &BaseItem{}
	err := pick.Create(engine, session, NewItemParams().
		withName("Rusty Pickaxe").
		withDescription("Rookie's pick, isn't capable of much.").withCallback(
		func(engine *akevitt.Akevitt, session *TemplateSession) error {
			AppendText(session, "Hello, world!", session.chat)

			return nil
		}).withQuantity(1))

	if err != nil {
		return err
	}

	character.Inventory = append(character.Inventory, pick)

	return character.Save(engine)
}

func (character *Character) Save(engine *akevitt.Akevitt) error {
	return engine.SaveGameObject(character, CharacterKey, character.account)
}

func (character *Character) GetDescription() string {
	return character.Description
}

func (character *Character) GetName() string {
	return character.Name
}

type CharacterParams struct {
	name        string
	description string
}
