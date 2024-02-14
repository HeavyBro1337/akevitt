package plugins

import (
	"time"

	"github.com/IvanKorchmit/akevitt"
)

type DefaultPlugins struct {
	Messages   *MessagePlugin
	HeartBeats *HeartBeatsPlugin
	plugins    []akevitt.Plugin
}

func (plugin *DefaultPlugins) Build(engine *akevitt.Akevitt) error {
	for _, p := range plugin.plugins {
		if err := p.Build(engine); err != nil {
			return err
		}
	}
	return nil
}

func NewDefaultPlugins() *DefaultPlugins {
	result := &DefaultPlugins{plugins: make([]akevitt.Plugin, 0)}

	result.Messages = NewMessagePlugin(true, nil, "")

	result.HeartBeats = NewHeartbeatPlugin().
		NewDuration(time.Second).
		NewDuration(time.Second * 10).
		NewDuration(time.Second * 30).
		NewDuration(time.Minute).
		Finish()

	result.plugins = append(result.plugins, result.Messages)
	result.plugins = append(result.plugins, result.HeartBeats)

	return result
}
