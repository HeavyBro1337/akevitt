package plugins

import (
	"errors"

	"github.com/IvanKorchmit/akevitt"
	"github.com/rivo/tview"
)

type DialogueFunc = func(engine *akevitt.Akevitt, session *akevitt.ActiveSession, dialogue *Dialogue) error

type DialoguePlugin struct {
	engine     *akevitt.Akevitt
	onDialogue DialogueFunc
}

// Dialogue struct for creating dialogue-like events.
// The content is actually a UI element, so you implement merchants, dialogues, etc.
// You must call NewDialogue to get started.
type Dialogue struct {
	title   string
	content tview.Primitive
	options []*Dialogue
	plugin  *DialoguePlugin
}

func NewDialoguePlugin(fn DialogueFunc) *DialoguePlugin {
	return &DialoguePlugin{
		onDialogue: fn,
	}
}

func (plugin *DialoguePlugin) Build(engine *akevitt.Akevitt) error {
	plugin.engine = engine
	return nil
}

// Creates new instance of dialogue
func (*DialoguePlugin) NewDialogue(title string) *Dialogue {
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
	res := dial.plugin.NewDialogue(title).SetContent(content)

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
func (dial *Dialogue) Proceed(index int, session *akevitt.ActiveSession, plugin *DialoguePlugin) error {
	return plugin.Dialogue(dial.options[index], session)
}

// Ends the dialogue with "Close" option
func (dial *Dialogue) End(text string) *Dialogue {
	dial.AddOption(text, nil)
	return dial
}

// Invokes dialogue event.
// Make sure you have installed the hook during initalisation.
func (plugin *DialoguePlugin) Dialogue(dialogue *Dialogue, session *akevitt.ActiveSession) error {
	if plugin.onDialogue == nil {
		return errors.New("dialogue callback is not installed")
	}

	return plugin.onDialogue(plugin.engine, session, dialogue)
}
