package akevitt

type Object interface {
	Save(engine *Akevitt) error // Save object into database
}

// Named object which has name and description
type NamedObject interface {
	GetName() string
	GetDescription() string
}

// Base of the gameobjects that can be involved in gameplay.
type GameObject interface {
	Object
	NamedObject
	Create(engine *Akevitt, session ActiveSession, params interface{}) error // Create an object wuth specified parameters.
}

// Container of game objects.
type Room interface {
	Object
	NamedObject
	GetExits() []Exit
	SetExits(exits ...Exit)
	GetKey() uint64
	GetObjects() []GameObject
	ContainObjects(objects ...GameObject)
	RemoveObject(object GameObject)
	OnCreate()
}

// The bridge between rooms
type Exit interface {
	Object
	GetRoom() Room
	GetKey() uint64
	SetRoom(room Room)
	Enter(engine *Akevitt, session ActiveSession) error
}
