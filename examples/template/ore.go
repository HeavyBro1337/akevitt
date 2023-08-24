package main

import (
	"akevitt"
	"fmt"
)

type Ore struct {
	BaseItem
	Price int
}

func (ore *Ore) Use(engine *akevitt.Akevitt, session *TemplateSession, other akevitt.GameObject) error {
	return ore.onUse(engine, session)
}
func (ore *Ore) GetDescription() string { return ore.Description + fmt.Sprintf(" (%d$)", ore.Price) }
