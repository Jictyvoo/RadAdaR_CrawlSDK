package datatypes

import "github.com/wrapped-owls/goremy-di/remy"

type Factory[T any] interface {
	New() (T, error)
	NewPtr() (*T, error)
}

type GenericFactory[T any] struct {
	injector    remy.Injector
	useKey      string
	constructor func() T
}

func NewConstructorFactory[T any](constructor func() T) Factory[T] {
	return GenericFactory[T]{constructor: constructor}
}

func NewInjectionFactory[T any](injector remy.Injector, useKey string) Factory[T] {
	return GenericFactory[T]{injector: injector, useKey: useKey}
}

func (fac GenericFactory[T]) New() (T, error) {
	if fac.constructor != nil {
		return fac.constructor(), nil
	}

	var dependencyKey []string
	if fac.useKey != "" {
		dependencyKey = []string{fac.useKey}
	}
	return remy.DoGet[T](fac.injector, dependencyKey...)
}

func (fac GenericFactory[T]) NewPtr() (*T, error) {
	result, err := fac.New()
	if err != nil {
		return nil, err
	}

	return &result, nil
}
