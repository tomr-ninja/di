package di

import "sync"

type ServiceConstructor[S any] func(*Registry) S

// Registry is a collection of services.
type (
	Registry struct {
		services map[string]any // map[string]serviceContainer
	}
	serviceContainer[S any] struct {
		service     S
		constructor ServiceConstructor[S]
		once        *sync.Once
	}
)

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]any),
	}
}

func GetService[S any](r *Registry, name string) S {
	c, ok := r.services[name]
	if !ok {
		panic("tried to get unregistered service")
	}

	container, ok := c.(serviceContainer[S])
	if !ok {
		panic("wrong container type")
	}

	container.once.Do(func() {
		container.service = container.constructor(r)
		r.services[name] = c
	})

	return container.service
}

func SetService[S any](r *Registry, name string, constructor ServiceConstructor[S]) {
	r.services[name] = serviceContainer[S]{
		constructor: constructor,
		once:        &sync.Once{},
	}
}
