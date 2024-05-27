package plugins

import (
	"time"

	"github.com/IvanKorchmit/akevitt"
)

func DefaultPlugins() []akevitt.Plugin {
	result := make([]akevitt.Plugin, 0)

	result = append(result, NewHeartbeatPlugin().
		NewDuration(time.Second).
		NewDuration(time.Second*10).
		NewDuration(time.Second*30).
		NewDuration(time.Minute).Finish(),
	)
	result = append(result, NewMessagePlugin(true, nil, ""))

	return result
}
