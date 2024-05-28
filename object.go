package akevitt

type Object interface {
	GetName() string
}

type Room struct {
	Name       string
	Exits      []*Exit
	Objects    []Object
	OnPreEnter func(*Akevitt, *ActiveSession) error
}

func (room *Room) GetKey() uint64 {
	return hash(room.Name)
}

type Exit struct {
	Name string
	Room *Room
}
