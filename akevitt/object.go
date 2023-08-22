package akevitt

type Object interface {
	Save(engine *Akevitt) error // Save object into database
}

type NamedObject interface {
	GetName() string
	GetDescription() string
}

type GameObject interface {
	Object
	NamedObject
	Create(engine *Akevitt, session ActiveSession, params interface{}) error
}

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

type Exit interface {
	Object
	GetRoom() Room
	GetKey() uint64
	SetRoom(room Room)
	Enter(engine *Akevitt, session ActiveSession) error
}
