package main

import "akevitt"

type autocomplete = func(entry string, engine *akevitt.Akevitt, session *TemplateSession) []string

var autocompletion map[string]autocomplete = make(map[string]autocomplete)

// Initialize autocompletion entries which can autocomplete with addiotnal arguments
// Example: `npc M`
// May suggest `npc Maxwell Jensen`
func initAutocompletion() {
	autocompletion["interact"] = func(entry string, engine *akevitt.Akevitt, session *TemplateSession) []string {
		npcs := akevitt.LookupOfType[*NPC](session.character.currentRoom)

		return akevitt.MapSlice(npcs, func(v *NPC) string {
			return "interact " + v.Name
		})
	}

	autocompletion["mine"] = func(entry string, engine *akevitt.Akevitt, session *TemplateSession) []string {
		ores := akevitt.LookupOfType[*Ore](session.character.currentRoom)

		return akevitt.MapSlice(ores, func(v *Ore) string {
			return "mine " + v.Name
		})
	}

	autocompletion["look"] = func(entry string, engine *akevitt.Akevitt, session *TemplateSession) []string {
		gameobjects := engine.Lookup(session.character.currentRoom)

		return akevitt.MapSlice(gameobjects, func(v akevitt.GameObject) string {
			return "look " + v.GetName()
		})
	}

}
