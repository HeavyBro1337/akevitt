package plugins

import (
	"fmt"
	"time"

	"github.com/IvanKorchmit/akevitt"
)

type HeartbeatMap = map[time.Duration]*akevitt.Pair[time.Ticker, []func() error]

type HeartBeatsPlugin struct {
	heartbeats HeartbeatMap
}

func (plugin *HeartBeatsPlugin) Build(engine *akevitt.Akevitt) error {
	for k := range plugin.heartbeats {
		plugin.startHeartBeats(k)
	}

	return nil
}

func NewHeartbeatPlugin() *HeartBuilder {
	return &HeartBuilder{
		plugin: &HeartBeatsPlugin{heartbeats: make(HeartbeatMap)},
	}
}

type HeartBuilder struct {
	plugin *HeartBeatsPlugin
}

func (builder *HeartBuilder) NewDuration(duration time.Duration) *HeartBuilder {
	builder.plugin.heartbeats[duration] = &akevitt.Pair[time.Ticker, []func() error]{L: *time.NewTicker(duration), R: make([]func() error, 0)}
	return builder
}

func (builder *HeartBuilder) Finish(duration time.Duration) *HeartBeatsPlugin {
	return builder.plugin
}

func (plugin *HeartBeatsPlugin) startHeartBeats(interval time.Duration) {
	go func() {
		t, ok := plugin.heartbeats[interval]
		errResults := make([]int, 0)
		if !ok {
			akevitt.LogWarn(fmt.Sprintf("ticker %d does not exist", interval))
			return
		}
		for range t.L.C {
			for i, fn := range t.R {
				if fn == nil {
					continue
				}
				if fn() != nil {
					errResults = append(errResults, i)
				}
			}

			for i := len(errResults) - 1; i >= 0; i-- {
				t.R = akevitt.RemoveItemByIndex(t.R, i)
			}
		}

	}()

}

func (plugin *HeartBeatsPlugin) SubscribeToHeartBeat(interval time.Duration, fn func() error) error {
	t, ok := plugin.heartbeats[interval]

	if !ok {
		return fmt.Errorf("warn: ticker %d does not exist", interval)
	}
	t.R = append(t.R, fn)
	return nil
}
