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

type ItemInterface interface {
	SetQuantity(value int)
	SetName(value string)
	SetDescription(value string)
	SetCallback(f InteractFunc)
}

func (item *Item) SetQuantity(value int) {
	item.Quantity = value
}

func (item *Item) SetName(value string) {
	item.Name = value
}

func (item *Item) SetDescription(value string) {
	item.Description = value
}

func (item *Item) SetCallback(f InteractFunc) {
	item.onUse = f
}

type Ore struct {
	Item
}

func (ore *Ore) Use(engine *akevitt.Akevitt, session *ActiveSession, other akevitt.GameObject) error {
	return ore.onUse(engine, session)
}

type HandItem struct {
	Item
}

type ItemParams struct {
	Quantity    int
	Name        string
	Description string
	onUse       InteractFunc
}

func NewItemParams() *ItemParams {
	return &ItemParams{Quantity: 1}
}

func (ip *ItemParams) withName(name string) *ItemParams {
	ip.Name = name
	return ip
}
func (ip *ItemParams) withDescription(description string) *ItemParams {
	ip.Description = description
	return ip
}
func (ip *ItemParams) withQuantity(quantity int) *ItemParams {
	ip.Quantity = quantity
	return ip
}
func (ip *ItemParams) withCallback(f InteractFunc) *ItemParams {
	ip.onUse = f
	return ip
}

func (item *Item) Create(engine *akevitt.Akevitt, session akevitt.ActiveSession, params interface{}) error {
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

func mapItemParams(item ItemInterface, itemParams *ItemParams) {
	item.SetName(itemParams.Name)
	item.SetDescription(itemParams.Description)
	item.SetCallback(itemParams.onUse)
	item.SetQuantity(itemParams.Quantity)
}

func createItem[T ItemInterface](item T, ip *ItemParams) T {
	mapItemParams(item, ip)

	return item

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

func (item *Item) Interact(engine *akevitt.Akevitt, session *ActiveSession) error {
	if item.onUse == nil {
		return errors.New("this item is unusable")
	}

	return item.onUse(engine, session)
}
