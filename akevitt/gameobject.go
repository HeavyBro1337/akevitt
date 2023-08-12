package akevitt

type Object interface {
	GetName() string
	Description() string                    // Retrieve description about that object
	Save(key uint64, engine *Akevitt) error // Save object into database
	OnLoad(engine *Akevitt) error
}

type GameObject interface {
	Object
	Create(engine *Akevitt, session *ActiveSession, params interface{}) error
	GetMap() map[string]Object
	OnRoomLookup() uint64
}
