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

func mapItemParams(item *Item, itemParams *ItemParams) {
	item.Name = itemParams.Name
	item.Description = itemParams.Description
	item.Quantity = itemParams.Quantity
}

func createItem(ip *ItemParams) *Item {
	item := &Item{}
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

func (item *Item) Interact(engine *akevitt.Akevitt, session akevitt.ActiveSession) error {
	sess, ok := session.(*ActiveSession)

	if !ok {
		return errors.New("could not cast to session")
	}

	if item.onUse == nil {
		return errors.New("this item is unusable")
	}

	return item.onUse(engine, sess)
}
