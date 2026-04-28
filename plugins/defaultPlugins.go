package plugins

import (
	"time"

	"github.com/IvanKorchmit/akevitt/internal/engine"
)

func DefaultPlugins() []engine.Plugin {
	result := make([]engine.Plugin, 0)

	result = append(result, NewHeartbeatPlugin().
		NewDuration(time.Second).
		NewDuration(time.Second*10).
		NewDuration(time.Second*30).
		NewDuration(time.Minute).Finish(),
	)
	result = append(result, NewMessagePlugin(true, nil, ""))

	return result
}