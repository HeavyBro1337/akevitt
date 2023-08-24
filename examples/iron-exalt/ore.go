package main

import (
	"fmt"

	"github.com/IvanKorchmit/akevitt"
)

type Ore struct {
	BaseItem
	Price int
}

func (ore *Ore) Use(engine *akevitt.Akevitt, session *IronExaltSession, other akevitt.GameObject) error {
	return ore.onUse(engine, session)
}
func (ore *Ore) GetDescription() string { return ore.Description + fmt.Sprintf(" (%d$)", ore.Price) }
