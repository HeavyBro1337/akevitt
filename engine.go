package akevitt

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"sync"
)

type Engine struct {
	plugins     []Plugin
	middlewares []ContextFunc
	wg          sync.WaitGroup
	rw          sync.RWMutex

	sessions map[*Context]bool

	handlers []Handler
}

func (e *Engine) AddSession(ctx *Context) {
	fmt.Printf("e.sessions: %v\n", e.sessions)
	e.rw.Lock()
	defer e.rw.Unlock()

	e.sessions[ctx] = true

	purgeDeadSessions(e)
}

func (e *Engine) RemoveSession(ctx *Context) {
	e.rw.Lock()
	defer e.rw.Unlock()

	delete(e.sessions, ctx)
}

func (e *Engine) Sessions() []*Context {
	contexts := make([]*Context, 0, len(e.sessions))
	for k := range maps.Keys(e.sessions) {
		contexts = append(contexts, k)
	}

	return contexts
}

func (e *Engine) Run() error {
	var errs error

	for _, p := range e.plugins {
		errs = errors.Join(p.Build(e.wrap()), errs)
	}

	for _, h := range e.handlers {
		e.wg.Add(1)
		log.Printf("[%s] Initializing...", h.Name())
		go func() {
			if err := h.Run(e); err != nil {
				log.Printf("[%s]: %v\n", h.Name(), err)
			}
			e.wg.Done()
		}()
	}
	e.wg.Wait()
	return errs
}

func NewGame(handlers ...Handler) *Builder {
	e := &Engine{
		plugins:     make([]Plugin, 0),
		handlers:    handlers,
		middlewares: make([]ContextFunc, 0),
		wg:          sync.WaitGroup{},
		rw:          sync.RWMutex{},
		sessions:    make(map[*Context]bool),
	}

	return &Builder{engine: e}
}

func (e *Engine) InvokeContext(ctx *Context) {
	for _, m := range e.middlewares {
		m(ctx)
	}
}

func (e *Engine) wrap() *Builder {
	return &Builder{
		engine: e,
	}
}
