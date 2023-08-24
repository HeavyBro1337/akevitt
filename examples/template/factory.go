package main

import (
	"akevitt"
)

func createOre(name, description string, price int) *Ore {
	oreParams := NewItemParams().
		withName(name).
		withDescription(description)
	oreParamsNoCallback := *oreParams
	oreParams.withCallback(func(engine *akevitt.Akevitt, session *TemplateSession) error {
		session.character.Inventory = append(session.character.Inventory, createItem(&Ore{Price: price}, &oreParamsNoCallback))
		return session.character.Save(engine)
	})
	ore := createItem(&Ore{Price: price}, oreParams)
	return ore
}

func createItem[T Item](item T, ip *ItemParams) T {
	mapItemParams(item, ip)

	return item
}

func mapItemParams(item Item, itemParams *ItemParams) {
	item.SetName(itemParams.Name)
	item.SetDescription(itemParams.Description)
	item.SetCallback(itemParams.onUse)
	item.SetQuantity(itemParams.Quantity)
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

func createNpc(name, description string) *NPC {
	return &NPC{Name: name, Description: description}
}

var lastRoomKey uint64 = 0

func createRoom(name, description string) *Room {
	r := &Room{Name: name, DescriptionData: description, Key: lastRoomKey}
	lastRoomKey++

	return r
}
