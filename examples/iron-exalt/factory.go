package main

import (
	"akevitt/akevitt"
)

func createOre(name, description string, price int) *Ore {
	oreParams := NewItemParams().
		withName(name).
		withDescription(description)
	oreParamsNoCallback := *oreParams
	oreParams.withCallback(func(engine *akevitt.Akevitt, session *ActiveSession) error {
		session.character.Inventory = append(session.character.Inventory, createItem(&Ore{Price: price}, &oreParamsNoCallback))
		return session.character.Save(engine)
	})
	ore := createItem(&Ore{Price: price}, oreParams)
	return ore
}
