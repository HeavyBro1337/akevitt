package akevitt

import (
	"golang.org/x/term"
)

type Plugin interface {
	Build(*Builder) error
}

type ShellPlugin struct {
	prompt   string
	callback func(*Context, *Engine, string)
}

func NewShellPlugin(prompt string, callback func(*Context, *Engine, string)) *ShellPlugin {
	return &ShellPlugin{
		prompt:   prompt,
		callback: callback,
	}
}

func (s *ShellPlugin) Build(b *Builder) error {
	b.Handle(func(ctx *Context) {
		t := term.NewTerminal(ctx, "> ")

		for {
			msg, err := t.ReadLine()

			if err != nil {
				b.engine.RemoveSession(ctx)
				return
			}

			s.callback(ctx, b.Engine(), msg)
		}
	})

	return nil
}
