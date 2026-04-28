package akevitt_test

import (
	"testing"

	"github.com/IvanKorchmit/akevitt"
	"github.com/stretchr/testify/assert"
)

func TestSplitCommand(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"look", []string{"look"}},
		{"look around", []string{"look", "around"}},
		{"go north", []string{"go", "north"}},
		{`say "hello world"`, []string{"say", "hello world"}},
		{"go north east", []string{"go", "north", "east"}},
		{"", nil},
	}

	for _, tc := range tests {
		result := akevitt.SplitCommand(tc.input)
		assert.Equal(t, tc.expected, result, "input: %s", tc.input)
	}
}

func TestLuaVM(t *testing.T) {
	vm := akevitt.NewLuaVM()
	assert.NotNil(t, vm)
	vm.Close()
}

func TestLuaVMReload(t *testing.T) {
	vm := akevitt.NewLuaVM()
	err := vm.Reload()
	assert.NoError(t, err)
	vm.Close()
}

func TestNewRoom(t *testing.T) {
	room := akevitt.NewRoom("Test Room")
	assert.Equal(t, "Test Room", room.Name)
	assert.NotEmpty(t, room.GUID)
	assert.NotNil(t, room.Exits)
	assert.NotNil(t, room.Objects)
}

func TestNewNPC(t *testing.T) {
	npc := akevitt.NewNPC("Guard")
	assert.Equal(t, "Guard", npc.Name)
	assert.NotEmpty(t, npc.GUID)
}

func TestNewItem(t *testing.T) {
	item := akevitt.NewItem("Sword")
	assert.Equal(t, "Sword", item.Name)
	assert.NotEmpty(t, item.GUID)
	assert.NotNil(t, item.Properties)
}
