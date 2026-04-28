package akevitt

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var objectCounter uint64

func generateGUID() string {
	id := atomic.AddUint64(&objectCounter, 1)
	return fmt.Sprintf("obj_%d", id)
}

type Object interface {
	GetName() string
	GetGUID() string
}

type ObjectImpl struct {
	Name        string
	Description string
	GUID        string
}

func (o *ObjectImpl) GetName() string {
	return o.Name
}

func (o *ObjectImpl) GetGUID() string {
	if o.GUID == "" {
		o.GUID = generateGUID()
	}
	return o.GUID
}

type Room struct {
	ObjectImpl
	Exits      []*Exit
	Objects    []Object
	OnPreEnter func(*Akevitt, *ActiveSession, *Exit) error
	mu         sync.RWMutex
}

func NewRoom(name string) *Room {
	return &Room{
		ObjectImpl: ObjectImpl{
			Name: name,
			GUID: generateGUID(),
		},
		Exits:   make([]*Exit, 0),
		Objects: make([]Object, 0),
	}
}

func (room *Room) AddObject(obj Object) {
	room.mu.Lock()
	defer room.mu.Unlock()
	room.Objects = append(room.Objects, obj)
}

func (room *Room) RemoveObject(obj Object) {
	room.mu.Lock()
	defer room.mu.Unlock()
	room.Objects = RemoveItem(room.Objects, obj)
}

func (room *Room) GetObjects() []Object {
	room.mu.RLock()
	defer room.mu.RUnlock()
	result := make([]Object, len(room.Objects))
	copy(result, room.Objects)
	return result
}

func (room *Room) Enter(engine *Akevitt, session *ActiveSession, targetExit *Exit) error {
	belongs := room == targetExit.Room

	if !belongs {
		return fmt.Errorf("the exit does not belong to %s", room.Name)
	}

	if room.OnPreEnter != nil {
		err := room.OnPreEnter(engine, session, targetExit)
		if err != nil {
			return err
		}
	}

	if targetExit.OnPreEnter != nil {
		err := targetExit.OnPreEnter(engine, session)
		if err != nil {
			return err
		}
	}
	return nil
}

func (room *Room) GetKey() uint64 {
	return hash(room.Name)
}

type Exit struct {
	Room       *Room
	Name       string
	OnPreEnter func(*Akevitt, *ActiveSession) error
}

type NPC struct {
	ObjectImpl
	RoomID    string
	Dialogue  string
	OnInteract func(*Akevitt, *ActiveSession) error
}

func NewNPC(name string) *NPC {
	return &NPC{
		ObjectImpl: ObjectImpl{
			Name: name,
			GUID:  generateGUID(),
		},
	}
}

type Item struct {
	ObjectImpl
	Properties map[string]any
}

func NewItem(name string) *Item {
	return &Item{
		ObjectImpl: ObjectImpl{
			Name: name,
			GUID:  generateGUID(),
		},
		Properties: make(map[string]any),
	}
}
