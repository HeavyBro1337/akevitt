package akevitt

import "fmt"

type Object interface {
	GetName() string
}

type Room struct {
	Name       string
	Exits      []*Exit
	Objects    []Object
	OnPreEnter func(*Akevitt, *ActiveSession, *Exit) error
}

func (room *Room) Enter(engine *Akevitt, session *ActiveSession, targetExit *Exit) error {
	belongs := Find(room.Exits, targetExit)

	if !belongs {
		return fmt.Errorf("the exit %s does not belong to %s", targetExit.Name, room.Name)
	}

	if room.OnPreEnter != nil {
		return room.OnPreEnter(engine, session, targetExit)
	}

	return nil
}

func (room *Room) GetKey() uint64 {
	return hash(room.Name)
}

type Exit struct {
	Name string
	Room *Room
}
