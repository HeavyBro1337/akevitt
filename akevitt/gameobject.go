package akevitt

type Object interface {
	Save(engine *Akevitt) error // Save object into database
}

type NamedObject interface {
	GetName() string
	Description() string
}

type GameObject interface {
	Object
	Create(engine *Akevitt, session ActiveSession, params interface{}) error
}
