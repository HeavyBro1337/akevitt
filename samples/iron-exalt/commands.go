package main

import "akevitt/akevitt"

func ooc(engine *akevitt.Akevitt, session *akevitt.ActiveSession, command string) {
	engine.SendOOCMessage(command, session)
}
