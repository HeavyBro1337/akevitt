package main

import (
	"akevitt/akevitt"
	"errors"
	"strconv"
)

func enter(engine *akevitt.Akevitt, session akevitt.ActiveSession, command string) error {
	sess, ok := session.(*ActiveSession)
	if !ok {
		return errors.New("invalid session type")
	}
	character := sess.character

	roomKey, err := strconv.ParseUint(command, 10, 64)
	if err != nil {
		return err
	}
	exit, err := akevitt.IsRoomReachable[*Room](engine, session, roomKey, character.CurrentRoomKey)
	if err != nil {
		return err
	}
	return exit.Enter(engine, session)
}
