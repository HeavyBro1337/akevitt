package akevitt

type Room interface {
	Object
	NamedObject
	GetExits() []Exit
	SetExits(exits ...Exit)
	GetKey() uint64
}

type Exit interface {
	Object
	NamedObject
	GetRoom() Room
	GetKey() uint64
	SetRoom(room Room)
	Enter(engine *Akevitt, session ActiveSession) error
}

func BindRooms[T Exit](room Room, sampleExit T, otherRooms ...Room) {
	var exits []Exit = make([]Exit, 0)
	for _, v := range otherRooms {
		if v == room {
			continue
		}
		exit := sampleExit
		exit.SetRoom(v)
		exits = append(exits, exit)
	}
	room.SetExits(exits...)
}
