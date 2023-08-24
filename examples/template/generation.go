package main

import (
	"akevitt"
	"fmt"

	"github.com/rivo/tview"
)

func generateRooms() *Room {
	spawn := createRoom("Spawn Room", "Just a spawn room")
	npcs := generateNpcs()

	mine := createRoom("Mine", "Mine of the corporation.")

	mine.AddObjects(
		createOre("Iron Ore", "Okay, this is epic", 6),
		createOre("Copper Ore", "Hmmmm, I wonder what happens if you mix it with tin?", 6),
		createOre("Tin Ore", "Hmmmm, I wonder what happens if you mix it with copper?", 6))
	spawn.AddObjects(npcs...)

	akevitt.BindRooms[*Exit](spawn, mine)
	akevitt.BindRooms[*Exit](mine, spawn)
	return spawn
}

func generateNpcs() []akevitt.GameObject {
	return []akevitt.GameObject{
		createNpc("Maxwell Jensen", "The tutor").
			UseInteract(func(engine *akevitt.Akevitt, session *TemplateSession) error {
				lore := tview.NewTextView().SetText(`
	Welcome, fellow miner! In our corporation you will mine minerals and ores for us!
	Go to the mine and mine with your pickaxe. 
	Come back and ask Ivan to deposite ores and receive the salary.
	Good luck!
	`)
				cya := tview.NewTextView().SetText(`See you later!`)
				d := akevitt.NewDialogue("Welcome")
				d.SetContent(lore)
				d.
					AddOption("Ok", cya).
					End()
				return engine.Dialogue(d, session)
			}),
		createNpc("Ivan Korchmit", "Depositor").
			UseInteract(func(engine *akevitt.Akevitt, session *TemplateSession) error {
				d := akevitt.NewDialogue("Deposit").SetContent(
					inventoryList[*Ore](engine, session,
						func(item Interactable) {
							AppendText(session, fmt.Sprintf("Sold %s", item.GetName()), session.chat)
							session.character.Inventory = akevitt.RemoveItem(session.character.Inventory, item)
							ore := item.(*Ore)
							session.character.Money += ore.Price
						})).End()

				return engine.Dialogue(d, session)
			})}
}
