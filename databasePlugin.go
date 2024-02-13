package akevitt

type DatabasePlugin[T Object] interface {
	Plugin

	Save(T) error
	LoadAll() ([]T, error)
}
