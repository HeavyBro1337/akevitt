package akevitt

// Accepts function which gets called when the user lefts the game.
// Note: use with caution, because calling methods from the engine like Message
// will cause an infinite recursion
// and in result: the application will crash.
func (builder *akevittBuilder) UseOnSessionEnd(f DeadSessionFunc) *akevittBuilder {
	if builder.engine.onDeadSession == nil {
		builder.engine.onDeadSession = make([]func(deadSession *ActiveSession, liveSessions []*ActiveSession, engine *Akevitt), 0)
	}
	builder.engine.onDeadSession = append(builder.engine.onDeadSession, f)

	return builder
}
