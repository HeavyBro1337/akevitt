package akevitt

type DatabasePlugin[T Object] interface {
	Plugin

	Save(Object) error
	LoadAll() ([]T, error)
}
