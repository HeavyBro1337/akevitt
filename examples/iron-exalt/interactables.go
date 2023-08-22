package main

import "akevitt"

type Interactable interface {
	akevitt.GameObject
	Interact(engine *akevitt.Akevitt, session *ActiveSession) error
}

type Usable interface {
	Interactable
	Use(engine *akevitt.Akevitt, session *ActiveSession, other akevitt.GameObject) error
}
