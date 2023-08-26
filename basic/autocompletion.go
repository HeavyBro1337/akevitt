package basic

import (
	"github.com/IvanKorchmit/akevitt"
)

type autocomplete = func(entry string, engine *akevitt.Akevitt, session *Session) []string

var autocompletion map[string]autocomplete = make(map[string]autocomplete)

// Initialize autocompletion entries which can autocomplete with addiotnal arguments
// Example: `npc M`
// May suggest `npc Maxwell Jensen`
func InitAutocompletion() {
	autocompletion["interact"] = func(entry string, engine *akevitt.Akevitt, session *Session) []string {
		npcs := akevitt.LookupOfType[*NPC](session.Character.currentRoom)

		return akevitt.MapSlice(npcs, func(v *NPC) string {
			return "interact " + v.Name
		})
	}

	autocompletion["look"] = func(entry string, engine *akevitt.Akevitt, session *Session) []string {
		gameobjects := session.Character.currentRoom.GetObjects()

		return akevitt.MapSlice(gameobjects, func(v akevitt.GameObject) string {
			return "look " + v.GetName()
		})
	}

}
