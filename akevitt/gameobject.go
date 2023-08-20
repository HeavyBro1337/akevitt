package akevitt

type Object interface {
	Save(engine *Akevitt) error // Save object into database
}

type NamedObject interface {
	GetName() string
	GetDescription() string
}

type GameObject interface {
	Object
	NamedObject
	Create(engine *Akevitt, session ActiveSession, params interface{}) error
}

type Interactable interface {
	GameObject
	Interact(engine *Akevitt, session ActiveSession) error
}

type Usable interface {
	Interactable
	Use(engine *Akevitt, session ActiveSession, other GameObject) error
}
