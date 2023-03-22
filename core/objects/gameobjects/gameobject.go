package gameobjects

import (
	"akevitt/core/objects"
	"akevitt/core/objects/credentials"
	"fmt"
)

// In-game object that you can interact within the game.
type GameObject interface {
	*objects.Object
	Create() error
}
type Room struct {
	Name string
}
type Character struct {
	Name      string
	Health    int
	MaxHealth int
	account   *credentials.Account
}

func (c Character) Description() string {
	return fmt.Sprintf("Name: %s\nHealth: %d/%d", c.Name, c.Health, c.MaxHealth)
}

func (c Character) Create() error {
	c.Health = 100
	c.MaxHealth = c.Health
	c.Name = "John Doe"
	return nil
}
