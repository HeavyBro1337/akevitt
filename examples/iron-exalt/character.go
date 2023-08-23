package main

import (
	"akevitt"
	"errors"
	"fmt"
)

type Character struct {
	Name           string
	Health         int
	Money          int
	MaxHealth      int
	account        *akevitt.Account
	currentRoom    akevitt.Room
	Inventory      []Interactable
	CurrentRoomKey uint64
}

func (character *Character) Create(engine *akevitt.Akevitt, session akevitt.ActiveSession, params interface{}) error {
	fmt.Println("Creating character...")

	characterParams, ok := params.(CharacterParams)
	if !ok {
		return errors.New("invalid params given")
	}
	sess, ok := session.(*IronExaltSession)

	if !ok {
		return errors.New("invalid session type")
	}

	character.Name = characterParams.name
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
		func(engine *akevitt.Akevitt, session *IronExaltSession) error {
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
