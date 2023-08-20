package main

import "akevitt/akevitt"

type NPC struct {
	Name        string
	Description string
}

func createNpc(name, description string, key uint64) *NPC {
	return &NPC{Name: name, Description: description}
}

func (npc *NPC) Create(engine *akevitt.Akevitt, session akevitt.ActiveSession, params interface{}) error {
	return nil
}

func (npc *NPC) GetName() string {
	return npc.Name
}

func (npc *NPC) GetDescription() string {
	return npc.Description
}

func (npc *NPC) Save(engine *akevitt.Akevitt) error {
	return engine.SaveObject(npc, NpcKey)
}
