package akevitt

import (
	"fmt"
	"reflect"
)

type Plugin interface {
	Build(*Akevitt) error
}

func FetchPlugin[T Plugin](engine *Akevitt) (*T, error) {
	for _, plugin := range engine.plugins {
		tPlugin, ok := plugin.(T)

		if ok {
			return &tPlugin, nil
		}
	}

	tType := reflect.TypeOf(new(T)).Name()

	return nil, fmt.Errorf("couldn't fetch the plugin of type %s", tType)
}
