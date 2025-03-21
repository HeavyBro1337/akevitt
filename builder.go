package akevitt

type ContextFunc func(*Context)

type Builder struct {
	engine *Engine
}

func (b *Builder) Engine() *Engine {
	return b.engine
}

func (b *Builder) Handle(handle ContextFunc) *Builder {
	b.engine.middlewares = append(b.engine.middlewares, handle)

	return b
}

func (b *Builder) Plugin(p Plugin) *Builder {
	b.engine.plugins = append(b.engine.plugins, p)

	return b
}
