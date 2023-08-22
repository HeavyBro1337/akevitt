package main

import "akevitt/akevitt"

func createOre(name, description string, price int) *Ore {
	oreParams := NewItemParams().
		withName("Tin Ore").
		withDescription("Hmmmm, I wonder what happens if you mix it with copper?")
	oreParamsNoCallback := *oreParams
	oreParams.withCallback(func(engine *akevitt.Akevitt, session *ActiveSession) error {
		session.character.Inventory = append(session.character.Inventory, createItem(&Ore{}, &oreParamsNoCallback))
		return session.character.Save(engine)
	})

	return createItem(&Ore{}, oreParams)
}
