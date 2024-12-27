package ui

func (d *DOM) OnCallback(key string, callback func()) {
	d.callbacks[key] = callback
}

func (d *DOM) OnStringEvent(key string, callback func(text string)) {
	d.stringCallbacks[key] = callback
}

func (d *DOM) OnBooleanEvent(key string, callback func(checked bool)) {
	d.boolCallbacks[key] = callback
}

func (d *DOM) OnIntEvent(key string, callback func(value int)) {
	d.intCallbacks[key] = callback
}
