package ironexalt

import (
	"akevitt"
	"errors"
	"strconv"
)

func enter(engine *akevitt.Akevitt, session akevitt.ActiveSession, arguments string) error {
	sess, ok := session.(*ActiveSession)
	if !ok {
		return errors.New("invalid session type")
	}
	character := sess.character
	prevRoom := character.currentRoom.GetName()
	roomKey, err := strconv.ParseUint(arguments, 10, 64)
	if err != nil {
		return err
	}
	exit, err := akevitt.IsRoomReachable[*Room](engine, session, roomKey, character.CurrentRoomKey)
	if err != nil {
		return err
	}
	err = exit.Enter(engine, session)

	if err != nil {
		return err
	}
	engine.Message(prevRoom, "left room", character.Name, session)
	engine.Message(character.currentRoom.GetName(), "entered room", character.Name, session)
	return nil
}
