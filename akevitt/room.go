package akevitt

import "reflect"

type Room interface {
	Object
	NamedObject
	GetExits() []Exit
	SetExits(exits ...Exit)
	GetKey() uint64
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
