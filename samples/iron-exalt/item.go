package main

import (
	"akevitt/akevitt"
	"errors"
)

type Item struct {
	Quantity    int
	ID          uint64
	Name        string
	Description string
	onUse       InteractFunc
}

type ItemParams struct {
	Quantity    int
	Name        string
	Description string
	onUse       InteractFunc
}

func (item *Item) Create(engine *akevitt.Akevitt, session akevitt.ActiveSession, params interface{}) error {
	itemParams, ok := params.(*ItemParams)

	if !ok {
		return errors.New("could not cast to item parameters")
	}
	item.Name = itemParams.Name
	item.Description = itemParams.Description
	item.Quantity = itemParams.Quantity

	key, err := engine.GenerateKey(item)

	if err != nil {
		return err
	}

	item.ID = key
	item.onUse = itemParams.onUse

	return item.Save(engine)
}

func (item *Item) Save(engine *akevitt.Akevitt) error {
	return engine.SaveObject(item, item.ID)
}

func (item *Item) GetDescription() string {
	return item.Description
}

func (item *Item) GetName() string {
	return item.Name
}
