package akevitt

import (
	"github.com/rivo/tview"
)

// Dialogue struct for creating dialogue-like events.
// The content is actually a UI element, so you implement merchants, dialogues, etc.
// You must call NewDialogue to get started.
type Dialogue struct {
	title   string
	content tview.Primitive
	options []*Dialogue
}

// Creates new instance of dialogue
func NewDialogue(title string) *Dialogue {
	dial := &Dialogue{title: title}
	dial.options = make([]*Dialogue, 0)

	return dial
}

// Accepts UI element to be displayed on user.
func (dial *Dialogue) SetContent(content tview.Primitive) *Dialogue {
	dial.content = content

	return dial
}

// Add options to the dialogue
func (dial *Dialogue) AddOption(title string, content tview.Primitive) *Dialogue {
	if dial.options == nil {
		dial.options = make([]*Dialogue, 0)
	}
	res := NewDialogue(title).SetContent(content)

	dial.options = append(dial.options, res)
	return res
}

// Gets an array of options
func (dial *Dialogue) GetOptions() []*Dialogue {
	return dial.options
}

// Gets the UI of dialogue
func (dial *Dialogue) GetContents() tview.Primitive {
	return dial.content
}

// Gets title of an option
func (dial *Dialogue) GetTitle() string {
	return dial.title
}

// Invokes engine.Dialogue of the specified index from options.
func (dial *Dialogue) Proceed(index int, session *ActiveSession, engine *Akevitt) error {
	return engine.Dialogue(dial.options[index], session)
}

// Ends the dialogue with "Close" option
func (dial *Dialogue) End() *Dialogue {
	dial.AddOption("Close", nil)
	return dial
}
