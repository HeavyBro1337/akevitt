package main

import (
	"akevitt/akevitt"
	"errors"
	"fmt"
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

func (room *Room) Create(engine *akevitt.Akevitt, session *akevitt.ActiveSession, params interface{}) error {
	return errors.New("room create is unused")
}

func (room *Room) GetName() string {
	return room.Name
}

func (room *Room) Description() string {
	return fmt.Sprintf("%s: %s", room.Name, room.DescriptionData)
}

func (room *Room) GetExits() []akevitt.Exit {
	return room.Exits
}

func (room *Room) SetExits(exits ...akevitt.Exit) {
	room.Exits = exits
}

func (room *Room) Save(key uint64, engine *akevitt.Akevitt) error {
	return errors.New("must not save room in database")
}

func (exit *Exit) Description() string {
	return "Hello!"
}

func (exit *Exit) Save(key uint64, engine *akevitt.Akevitt) error {
	return engine.SaveWorldObject(exit, key)
}

type Exit struct {
	Key  uint64
	room akevitt.Room
}

func (exit *Exit) GetName() string {
	if exit.room == nil {
		return "Nowhere"
	}

	return exit.room.GetName()
}

func (exit *Exit) Enter(engine *akevitt.Akevitt, session *akevitt.ActiveSession) error {
	character, ok := session.RelatedGameObjects[currentCharacterKey].Second.(*Character)
	if !ok {
		return errors.New("could not cast to character")
	}

	character.CurrentRoomKey = exit.Key
	actualRoom, err := akevitt.GetStaticObject[akevitt.Room](engine, exit.Key)

	if err != nil {
		return err
	}

	character.currentRoom = actualRoom

	engine.SendRoomMessage("Entered room", session)

	return character.Save(session.RelatedGameObjects[currentCharacterKey].First, engine)
}

func (exit *Exit) GetRoom() akevitt.Room {
	return exit.room
}

func (room *Room) OnLoad(engine *akevitt.Akevitt) error {
	println("Invoked OnLoad in room")
	for _, v := range room.Exits {
		otherRoom, err := akevitt.GetStaticObject[akevitt.Room](engine, v.GetKey())

		if err != nil {
			return err
		}

		v.SetRoom(otherRoom)
	}
	return nil
}

func (exit *Exit) OnLoad(engine *akevitt.Akevitt) error {
	return nil
}

func (exit *Exit) GetKey() uint64 {
	return exit.Key
}

func (exit *Exit) SetRoom(room akevitt.Room) {
	exit.Key = room.GetKey()
	exit.room = room
}

func BindRooms[T akevitt.Exit](room akevitt.Room, sampleExit T, otherRooms ...akevitt.Room) {
	var exits []akevitt.Exit = make([]akevitt.Exit, 0)

	for _, v := range otherRooms {
		if v == room {
			continue
		}
		exit := sampleExit

		fmt.Printf("v.GetName(): %v\n", v.GetName())
		fmt.Printf("room.GetName(): %v\n", room.GetName())

		exit.SetRoom(v)
		exits = append(exits, exit)
	}

	room.SetExits(exits...)
}
