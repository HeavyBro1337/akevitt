package akevitt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Hashes password using Bcrypt algorithm
func hashString(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

// Compares hash and password.
func compareHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Finds `T` value of []T.
func Find[T comparable](collection []T, value T) bool {
	for _, b := range collection {
		if b == value {
			return true
		}
	}
	return false
}

// Converts `Uint64` to byte array.
func intToByte(value uint64) []byte {
	binaryId := make([]byte, 8)
	binary.BigEndian.PutUint64(binaryId, uint64(value))
	return binaryId
}

func FindByKey[TCollection, T comparable](collection []TCollection, selector func(key TCollection) T, value T) *TCollection {
	if collection == nil {
		panic(errors.New("collection is nil"))
	}
	for _, b := range collection {
		if selector(b) == value {
			return &b
		}
	}
	return nil
}

// Removes item from collection and returns it.
func RemoveItem[T comparable](l []T, item T) []T {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func RemoveItemByIndex[T any](l []T, i int) []T {
	return append(l[:i], l[i+1:]...)
}

// Maps slice, similar to JavaScript's map method.
func MapSlice[T any, TResult any](l []T, callback func(v T) TResult) []TResult {
	result := make([]TResult, 0)

	for _, v := range l {
		result = append(result, callback(v))
	}

	return result
}

func FindNeighboringRoomByName(currentRoom *Room, name string) (*Room, *Exit, error) {
	for _, v := range currentRoom.Exits {
		if strings.EqualFold(v.Room.Name, name) {
			return v.Room, v, nil
		}
	}
	return nil, nil, fmt.Errorf("room %s not found", name)
}

// Checks if current room specified reachable to another room.
func IsRoomReachable[T Room](engine *Akevitt, session *ActiveSession, name string, currentRoomKey uint64) (*Exit, error) {
	room, err := engine.GetRoom(currentRoomKey)

	if err != nil {
		return nil, err
	}

	exits := room.Exits

	if exits == nil {
		return nil, errors.New("array of exits is nil")
	}
	exit := FindByKey(exits, func(key *Exit) string {
		return strings.ToLower(key.Room.Name)
	}, strings.ToLower(name))
	if exit == nil {
		return nil, errors.New("unreachable")
	}
	return *exit, nil
}

// Binds room with an exit.
func BindRooms(emptyExit Exit, room *Room, otherRooms ...*Room) {
	exits := make([]*Exit, 0)
	for _, v := range otherRooms {
		exit := emptyExit
		exit.Room = v // Setting exit's current room
		exits = append(exits, &exit)
	}

	room.Exits = exits
}

// Saves game object in a database associated with an account.
func (engine *Akevitt) SaveGameObject(gameObject GameObject, key uint64, account *Account) error {
	if account == nil {
		return errors.New("account is nil")
	}

	databasePlugin, err := FetchPlugin[DatabasePlugin[GameObject]](engine)

	if err != nil {
		return err
	}
	return (*databasePlugin).Save(gameObject)
}

// Saves object into a database.
func (engine *Akevitt) SaveObject(gameObject GameObject) error {
	databasePlugin, err := FetchPlugin[DatabasePlugin[GameObject]](engine)

	if err != nil {
		return err
	}
	return (*databasePlugin).Save(gameObject)
}

func CreateObject[T GameObject](engine *Akevitt, session *ActiveSession, object T, params interface{}) (T, error) {
	return object, object.Create(engine, session, params)
}

func (engine *Akevitt) GlobalLookup(room *Room, name string) []GameObject {
	return globalSearchRecursive(engine.defaultRoom, name, nil, nil)
}

func LookupOfType[T GameObject](room Room) []T {
	return FilterByType[T, GameObject](room.Objects)
}

func globalSearchRecursive(room *Room, name string, visited []string, result []GameObject) []GameObject {
	if visited == nil {
		visited = make([]string, 0)
	}

	if room == nil {
		return nil
	}
	if result == nil {
		result = make([]GameObject, 0)
	}

	visited = append(visited, room.Name)

	for _, v := range room.Objects {
		if strings.EqualFold(v.GetName(), name) {
			result = append(result, v)
		}
	}

	for _, v := range room.Exits {
		r := v.Room

		if Find[string](visited, r.Name) {
			continue
		}

		err := globalSearchRecursive(r, name, visited, result)

		if err != nil {
			return err
		}
	}
	return nil
}

func FilterByType[T any, TCollection any](collection []TCollection) []T {
	result := make([]T, 0)
	for _, v := range collection {
		t, ok := any(v).(T)

		if ok {
			result = append(result, t)
		}
	}
	return result
}

func saveRoomsRecursively(engine *Akevitt, room *Room, visited []string) error {
	if visited == nil {
		visited = make([]string, 0)
	}

	if room == nil {
		return errors.New("room is nil")
	}

	fmt.Printf("Loading Room: %s\n", room.Name)

	engine.rooms[room.GetKey()] = room

	visited = append(visited, room.Name)

	for _, v := range room.Exits {
		r := v.Room

		if Find[string](visited, r.Name) {
			continue
		}

		err := saveRoomsRecursively(engine, r, visited)

		if err != nil {
			return err
		}
	}
	return nil
}

func hash(s string) uint64 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return uint64(h.Sum32())
}
