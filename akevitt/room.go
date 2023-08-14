package akevitt

import (
	"errors"
	"reflect"
)

type Room interface {
	Object
	NamedObject
	GetExits() []Exit
	SetExits(exits ...Exit)
	GetKey() uint64
	GetObjects() []GameObject
	OnCreate()
}

type Exit interface {
	Object
	GetRoom() Room
	GetKey() uint64
	SetRoom(room Room)
	Enter(engine *Akevitt, session ActiveSession) error
}

func BindRooms[T Exit](room Room, otherRooms ...Room) {
	var emptyExit T
	var exits []Exit = make([]Exit, 0)
	for _, v := range otherRooms {
		if v == room {
			continue
		}
		exit := reflect.New(reflect.TypeOf(emptyExit).Elem()).Interface().(T)
		exit.SetRoom(v)
		exits = append(exits, exit)
	}
	room.SetExits(exits...)
}

func IsRoomReachable[T Room](engine *Akevitt, session ActiveSession, roomKey uint64, currentRoomKey uint64) (Exit, error) {
	room, err := engine.GetRoom(currentRoomKey)

	if err != nil {
		return nil, err
	}

	exits := room.GetExits()

	if exits == nil {
		return nil, errors.New("array of exits is nil")
	}

	exit := findByKey(exits, func(key Exit) uint64 {
		return key.GetKey()
	}, currentRoomKey)
	if exit == nil {
		return nil, errors.New("unreachable")
	}
	return *exit, nil
}
