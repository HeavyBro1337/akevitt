package main

import "akevitt/akevitt"

type autocomplete = func(entry string, engine *akevitt.Akevitt, session *ActiveSession) []string

var autocompletion map[string]autocomplete = make(map[string]autocomplete)
