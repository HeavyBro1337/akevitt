package main

import "akevitt/akevitt"

func ooc(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) {
	engine.SendOOCMessage(command, session)

}

func characterMessage(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) {
	engine.SendRoomMessage(command, session)
}
