package akevitt

type Room interface {
	Object
	GetExits() []Exit
	SetExits(exits ...Exit)
	GetKey() uint64
}

type Exit interface {
	Object
	GetRoom() Room
	GetKey() uint64
	SetRoom(room Room)
	Enter(engine *Akevitt, session *ActiveSession) error
}
