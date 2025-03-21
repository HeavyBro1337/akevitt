package akevitt

import "io"

func purgeDeadSessions(engine *Engine) {
	for ctx := range engine.sessions {
		_, err := io.WriteString(ctx, "\000")
		if err != nil {
			engine.RemoveSession(ctx)
			continue
		}
	}
}
