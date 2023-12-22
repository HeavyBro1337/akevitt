package akevitt

// Out-of-character chat command
func OocCmd(engine *Akevitt, session *ActiveSession, command string) error {
	return engine.Message("ooc", command, session.Account.Username, session)
}
