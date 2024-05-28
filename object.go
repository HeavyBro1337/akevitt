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
		return fmt.Errorf("the exit does not belong to %s", room.Name)
	}

	if room.OnPreEnter != nil {
		err := room.OnPreEnter(engine, session, targetExit)
		if err != nil {
			return err
		}
	}

	if targetExit.OnPreEnter != nil {
		err := targetExit.OnPreEnter(engine, session)
		if err != nil {
			return err
		}
	}
	return nil
}

func (room *Room) GetKey() uint64 {
	return hash(room.Name)
}

type Exit struct {
	Room       *Room
	OnPreEnter func(*Akevitt, *ActiveSession) error
}
