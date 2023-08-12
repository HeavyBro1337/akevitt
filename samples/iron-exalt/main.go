package main

import (
	"akevitt/akevitt"
	"log"
)

func main() {
	log.Fatal(akevitt.NewEngine().
		UseMouse().
		UseDBPath("data/iron-exalt.db").
		UseRootUI(rootScreen).
		Run())
}
