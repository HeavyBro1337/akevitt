package akevitt

import (
	"github.com/rivo/tview"
)

type Dialogue struct {
	title   string
	content tview.Primitive
	options []*Dialogue
}

func NewDialogue(title string) *Dialogue {
	dial := &Dialogue{title: title}
	dial.options = make([]*Dialogue, 0)

	return dial
}

func (dial *Dialogue) SetContent(content tview.Primitive) *Dialogue {
	dial.content = content

	return dial
}

func (dial *Dialogue) AddOption(title string, content tview.Primitive) *Dialogue {
	if dial.options == nil {
		dial.options = make([]*Dialogue, 0)
	}
	res := NewDialogue(title).SetContent(content)

	dial.options = append(dial.options, res)
	return res
}

func (dial *Dialogue) GetOptions() []*Dialogue {
	return dial.options
}

func (dial *Dialogue) GetContents() tview.Primitive {
	return dial.content
}

func (dial *Dialogue) GetTitle() string {
	return dial.title
}

func (dial *Dialogue) Proceed(index int, session ActiveSession, engine *Akevitt) error {
	return engine.Dialogue(dial.options[index], session)
}

func (dial *Dialogue) End() *Dialogue {
	dial.AddOption("Close", nil)
	return dial
}
