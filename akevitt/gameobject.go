package akevitt

type Object interface {
	Save(engine *Akevitt) error // Save object into database
}

type GameObject interface {
	Object
	Create(engine *Akevitt, session *ActiveSession, params interface{}) error
	GetMap() map[string]Object
	OnRoomLookup() uint64
}
