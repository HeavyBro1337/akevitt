package akevitt

type akevittBuilder struct {
	engine *Akevitt
}

func (builder *akevittBuilder) Finish() *Akevitt {
	return builder.engine
}
