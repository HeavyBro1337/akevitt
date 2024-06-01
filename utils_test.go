package akevitt_test

import (
	"fmt"
	"testing"

	"github.com/IvanKorchmit/akevitt"
	"github.com/stretchr/testify/assert"
)

func TestBindRooms(t *testing.T) {
	assrt := assert.New(t)

	spawn := akevitt.Room{
		Name: "Spawn",
	}

	otherRoom := akevitt.Room{
		Name: "Markte",
	}

	akevitt.BindRooms(&spawn, &otherRoom)

	roomsOfSpawn := akevitt.MapSlice(spawn.Exits, func(v *akevitt.Exit) *akevitt.Room {
		return v.Room
	})

	roomsOfOther := akevitt.MapSlice(otherRoom.Exits, func(v *akevitt.Exit) *akevitt.Room {
		return v.Room
	})

	assrt.True(akevitt.Find(roomsOfSpawn, &otherRoom), roomsOfSpawn)
	assrt.False(akevitt.Find(roomsOfOther, &spawn), roomsOfOther)
}

func TestMap(t *testing.T) {
	assrt := assert.New(t)

	input := []int{1, 2, 3, 4}

	output := akevitt.MapSlice(input, func(v int) string {
		return fmt.Sprintf("%d.0", v)

	})

	assrt.Equal([]string{"1.0", "2.0", "3.0", "4.0"}, output)
}
