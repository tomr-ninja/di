package indi

import (
	"sync"
)

type ServiceConstructor[S any] func(*Registry) S

// Registry is a collection of services.
type (
	Registry struct {
		services map[string]any // map[string]serviceDef
	}
	serviceDef[S any] struct {
		service     S
		constructor ServiceConstructor[S]
		once        sync.Once
	}
)

// NewRegistry creates a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		services: make(map[string]any),
	}
}

func SetService[S any](
	r *Registry, name string, constructor ServiceConstructor[S],
) {
	r.services[name] = &serviceDef[S]{
		constructor: constructor,
		once:        sync.Once{},
	}
}

func GetService[S any](r *Registry, name string) S {
	c, ok := r.services[name]
	if !ok {
		panic("tried to get unregistered service")
	}

	def, ok := c.(*serviceDef[S])
	if !ok {
		panic("wrong def type")
	}

	def.init(r)

	return def.service
}

func (def *serviceDef[S]) init(r *Registry) {
	def.once.Do(func() {
		def.service = def.constructor(r)
	})
}

func InitAll(r *Registry) {
	type initable interface {
		init(r *Registry)
	}

	var wg sync.WaitGroup
	wg.Add(len(r.services))

	for _, c := range r.services {
		def := c.(initable)

		go func(i initable) {
			i.init(r)
			wg.Done()
		}(def)
	}

	wg.Wait()
}
