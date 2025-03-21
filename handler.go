package akevitt

type Handler interface {
	Run(*Engine) error
	Name() string
}
