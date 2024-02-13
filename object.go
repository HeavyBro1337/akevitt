package akevitt

type Object interface {
	Save(engine *Akevitt) error // Save object into a database
}

// Base of the gameobjects that can be involved in gameplay.
type GameObject interface {
	Object
	GetName() string
	Create(engine *Akevitt, session *ActiveSession, params any) error // Create an object wuth specified parameters.
}

type Room struct {
	Name       string
	Exits      []*Exit
	Objects    []GameObject
	OnPreEnter func(*ActiveSession) error
}

func (room *Room) GetKey() uint64 {
	return hash(room.Name)
}

type Exit struct {
	Name string
	Room *Room
}
