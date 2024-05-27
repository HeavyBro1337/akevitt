package main

import (
	"github.com/IvanKorchmit/akevitt"
	"github.com/IvanKorchmit/akevitt/plugins"
)

func main() {
	akevitt.NewEngine().AddPlugin(plugins.DefaultPlugins()...)
}
