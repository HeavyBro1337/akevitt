package akevitt

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"

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

// Converts `T` to byte array.
func serialize[T Object](v T) ([]byte, error) {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// Converts byte array to T struct.
func deserialize[T Object](b []byte) (T, error) {
	var result T
	var decodeBuffer bytes.Buffer
	decodeBuffer.Write(b)
	dec := gob.NewDecoder(&decodeBuffer)
	err := dec.Decode(&result)
	if err != nil {
		return result, err
	}
	return result, err
}

func findByKey[TCollection, T comparable](collection []TCollection, selector func(key TCollection) T, value T) *TCollection {
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

// Maps slice, similar to JavaScript's map method.
func MapSlice[T any, TResult any](l []T, callback func(v T) TResult) []TResult {
	result := make([]TResult, 0)

	for _, v := range l {
		result = append(result, callback(v))
	}

	return result
}

// Checks if current room specified reachable to another room.
func IsRoomReachable[T Room](engine *Akevitt, session ActiveSession, roomKey uint64, currentRoomKey uint64) (Exit, error) {
	room, err := engine.GetRoom(currentRoomKey)

	if err != nil {
		return nil, err
	}

	exits := room.GetExits()

	if exits == nil {
		return nil, errors.New("array of exits is nil")
	}
	exit := findByKey(exits, func(key Exit) uint64 {
		return key.GetKey()
	}, roomKey)
	if exit == nil {
		return nil, errors.New("unreachable")
	}
	return *exit, nil
}

// Binds room with an exit on both sides.
func BindRooms[T Exit](room Room, otherRooms ...Room) {
	var emptyExit T
	var exits []Exit = make([]Exit, 0)
	for _, v := range otherRooms {
		if v == room {
			continue
		}
		exit := reflect.New(reflect.TypeOf(emptyExit).Elem()).Interface().(T)
		exit.SetRoom(v)
		exits = append(exits, exit)
	}
	room.SetExits(exits...)
}

// Saves object to database.
func SaveObject[T Object](engine *Akevitt, obj T, category string, key uint64) error {
	return overwriteObject[T](engine.db, key, category, obj)
}

// Finds game object associated with an account in database.
func FindObject[T GameObject](engine *Akevitt, session ActiveSession, key uint64) (T, error) {
	return findObject[T](engine.db, *session.GetAccount(), key)
}

// Saves game object in a database associated with an account.
func (engine *Akevitt) SaveGameObject(gameObject GameObject, key uint64, account *Account) error {
	return overwriteObject(engine.db, key, account.Username, gameObject)
}

// Saves object into a database.
func (engine *Akevitt) SaveObject(gameObject GameObject, key uint64) error {
	return overwriteObject(engine.db, key, gameObject.GetName(), gameObject)
}

// Auto-increment uint64 key by object's name.
func (engine *Akevitt) GenerateKey(gameobject GameObject) (uint64, error) {
	return generateKey(engine.db, gameobject.GetName())
}

func CreateObject[T GameObject](engine *Akevitt, session ActiveSession, object T, params interface{}) (T, error) {
	return object, object.Create(engine, session, params)
}

func (engine *Akevitt) Lookup(room Room) []GameObject {
	return room.GetObjects()
}

func LookupOfType[T GameObject](room Room) []T {
	return FilterByType[T, GameObject](room.GetObjects())
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

func saveRoomsRecursively(engine *Akevitt, room Room, visited []string) error {
	if visited == nil {
		visited = make([]string, 0)
	}

	if room == nil {
		return errors.New("room is nil")
	}

	fmt.Printf("Loading Room: %s\n", room.GetName())

	engine.rooms[room.GetKey()] = room

	visited = append(visited, room.GetName())

	for _, v := range room.GetExits() {
		r := v.GetRoom()

		if Find[string](visited, r.GetName()) {
			continue
		}

		err := saveRoomsRecursively(engine, r, visited)

		if err != nil {
			return err
		}
	}
	return nil
}
