package basic

import (
	"errors"
	"fmt"

	"github.com/IvanKorchmit/akevitt"
)

const (
	CharacterKey uint64 = iota + 1
	NpcKey
)

type Interactable interface {
	akevitt.GameObject
	Interact(engine *akevitt.Akevitt, session *Session) error
}

type Usable[T akevitt.ActiveSession] interface {
	Interactable
	Use(engine *akevitt.Akevitt, session *Session, other akevitt.GameObject) error
}

// The character struct represents an in-game character that the user will play as.
type Character struct {
	Name           string
	Description    string
	Health         int
	Money          int
	MaxHealth      int
	currentRoom    akevitt.Room
	account        *akevitt.Account
	Inventory      []Interactable
	CurrentRoomKey uint64
}

func (char *Character) SetCurrentRoom(room akevitt.Room) {
	char.currentRoom = room
}

func (char *Character) GetCurrentRoom() akevitt.Room {
	return char.currentRoom
}

func (character *Character) Create(engine *akevitt.Akevitt, session akevitt.ActiveSession, params interface{}) error {
	sess := CastSession[*Session](session)

	fmt.Println("Creating a character...")

	characterParams, ok := params.(CharacterParams)
	if !ok {
		return errors.New("invalid params given")
	}

	character.Name = characterParams.name
	character.Description = characterParams.description
	character.Health = 10
	character.MaxHealth = 10
	character.Money = 100
	character.currentRoom = engine.GetSpawnRoom()
	character.Inventory = make([]Interactable, 0)
	sess.Character = character
	character.account = sess.account
	room := engine.GetSpawnRoom()
	room.AddObjects(sess.Character)
	character.currentRoom = room
	character.CurrentRoomKey = character.currentRoom.GetKey()

	pick := &BaseItem{}
	err := pick.Create(engine, session, NewItemParams().
		WithName("Rusty Pickaxe").
		WithDescription("Rookie's pick, isn't capable of much.").WithCallback(
		func(engine *akevitt.Akevitt, session *Session) error {
			AppendText("Hello, world!", session.Chat)

			return nil
		}).WithQuantity(1))

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
