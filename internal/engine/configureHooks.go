package engine

func (builder *akevittBuilder) UseOnSessionEnd(f DeadSessionFunc) *akevittBuilder {
	if builder.engine.onDeadSession == nil {
		builder.engine.onDeadSession = make([]func(deadSession *ActiveSession, liveSessions []*ActiveSession, engine *Akevitt), 0)
	}
	builder.engine.onDeadSession = append(builder.engine.onDeadSession, f)

	return builder
}