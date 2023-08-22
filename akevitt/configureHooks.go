package akevitt

func (engine *Akevitt) UseMessage(f MessageFunc) *Akevitt {
	engine.onMessage = f

	return engine
}

func (engine *Akevitt) UseDialogue(f DialogueFunc) *Akevitt {
	engine.onDialogue = f

	return engine
}

// Provide some callback if session is ended. Note: Some methods are dangerious to call i.e. engine.Message,
// because it may invoke dead session cleanup which will cause stack overflow error and crash the application.
func (engine *Akevitt) UseOnSessionEnd(f DeadSessionFunc) *Akevitt {
	engine.onDeadSession = f
	return engine
}
