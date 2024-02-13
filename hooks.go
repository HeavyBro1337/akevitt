package akevitt

import (
	"errors"
)

// Invokes dialogue event.
// Make sure you have installed the hook during initalisation.
func (engine *Akevitt) Dialogue(dialogue *Dialogue, session *ActiveSession) error {
	if engine.onDialogue == nil {
		return errors.New("dialogue callback is not installed")
	}

	return engine.onDialogue(engine, session, dialogue)
}
