package akevitt

// Accepts the function which gets invoked when someone sends the message (engine.Message)
func (engine *Akevitt) UseOnMessage(f MessageFunc) *Akevitt {
	engine.onMessage = f

	return engine
}

// Called when engine.Dialogue is called
func (engine *Akevitt) UseOnDialogue(f DialogueFunc) *Akevitt {
	engine.onDialogue = f

	return engine
}

// Accepts function which gets called when the user lefts the game.
// Note: use with caution, because calling methods from the engine like Message
// will cause an infinite recursion
// and in result: the application will crash.
func (engine *Akevitt) UseOnSessionEnd(f DeadSessionFunc) *Akevitt {
	engine.onDeadSession = f
	return engine
}
