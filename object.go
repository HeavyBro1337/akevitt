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
	belongs := func() bool {
		for _, v := range targetExit.Room.Exits {
			if v.Room == room {
				return true
			}
		}

		return false
	}

	if !belongs() {
		return fmt.Errorf("the exit does not belong to %s", room.Name)
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
	Room *Room
}
