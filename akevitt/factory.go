package akevitt

import "errors"

type Factory[T GameObject] struct {
	result []T
	params []interface{}
}

func NewFactory[T GameObject]() *Factory[T] {
	factory := &Factory[T]{}
	factory.result = make([]T, 0)
	factory.params = make([]interface{}, 0)
	return factory
}

func (factory *Factory[T]) AddObject(obj T, params interface{}) *Factory[T] {
	factory.result = append(factory.result, obj)
	factory.params = append(factory.params, params)

	return factory
}

func (factory *Factory[T]) Finish(engine *Akevitt, session ActiveSession) ([]T, error) {
	errs := make([]error, 0)

	for k, v := range factory.result {
		errs = append(errs, v.Create(engine, session, factory.params[k]))
	}

	return factory.result, errors.Join(errs...)
}
