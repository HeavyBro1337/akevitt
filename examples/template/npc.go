package main

import (
	"akevitt"
	"fmt"
)

type NPC struct {
	Name        string
	Description string
	onInteract  InteractFunc
}

type InteractFunc = func(engine *akevitt.Akevitt, session *TemplateSession) error

func (npc *NPC) UseInteract(f InteractFunc) *NPC {
	npc.onInteract = f

	return npc
}

func (npc *NPC) Create(engine *akevitt.Akevitt, session akevitt.ActiveSession, params interface{}) error {
	return nil
}

func (npc *NPC) GetName() string        { return npc.Name }
func (npc *NPC) GetDescription() string { return npc.Description }

func (npc *NPC) Save(engine *akevitt.Akevitt) error {
	return engine.SaveObject(npc, NpcKey)
}

func (npc *NPC) Interact(engine *akevitt.Akevitt, session *TemplateSession) error {
	if npc.onInteract == nil {
		return fmt.Errorf("npc named %s has no interact callback installed", npc.Name)
	}
	return npc.onInteract(engine, session)
}
