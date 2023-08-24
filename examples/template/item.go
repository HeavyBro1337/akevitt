package main

import (
	"akevitt"
	"errors"
)

type Interactable interface {
	akevitt.GameObject
	Interact(engine *akevitt.Akevitt, session *TemplateSession) error
}

type Usable interface {
	Interactable
	Use(engine *akevitt.Akevitt, session *TemplateSession, other akevitt.GameObject) error
}

type BaseItem struct {
	Quantity    int
	ID          uint64
	Name        string
	Description string
	onUse       InteractFunc
}

type Item interface {
	akevitt.GameObject

	SetQuantity(value int)
	SetName(value string)
	SetDescription(value string)
	SetCallback(f InteractFunc)
}

func (item *BaseItem) SetQuantity(value int)       { item.Quantity = value }
func (item *BaseItem) SetName(value string)        { item.Name = value }
func (item *BaseItem) SetDescription(value string) { item.Description = value }
func (item *BaseItem) SetCallback(f InteractFunc)  { item.onUse = f }

type HandItem struct {
	BaseItem
}

func (item *BaseItem) Create(engine *akevitt.Akevitt, session akevitt.ActiveSession, params interface{}) error {
	itemParams, ok := params.(*ItemParams)

	if !ok {
		return errors.New("could not cast to item parameters")
	}
	mapItemParams(item, itemParams)

	key, err := engine.GenerateKey(item)

	if err != nil {
		return err
	}

	item.ID = key
	item.onUse = itemParams.onUse

	return item.Save(engine)
}

func (item *BaseItem) Save(engine *akevitt.Akevitt) error { return engine.SaveObject(item, item.ID) }
func (item *BaseItem) GetName() string                    { return item.Name }
func (item *BaseItem) GetDescription() string             { return item.Description }

func (item *BaseItem) Interact(engine *akevitt.Akevitt, session *TemplateSession) error {
	if item.onUse == nil {
		return errors.New("this item is unusable")
	}

	return item.onUse(engine, session)
}
