package akevitt

type Object interface {
	GetName() string
	Description() string                    // Retrieve description about that object
	Save(key uint64, engine *akevitt) error // Save object into database
	OnLoad(engine *akevitt) error
}

type GameObject interface {
	Object
	Create(engine *akevitt, session *ActiveSession, params interface{}) error
	GetMap() map[string]Object
	OnRoomLookup() uint64
}
