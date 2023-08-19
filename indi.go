package indi

import (
	"sync"
)

var defaultRegistry = NewRegistry()

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

func SetServiceFromRegistry[S any](r *Registry, name string, constructor ServiceConstructor[S]) {
	r.services[name] = &serviceDef[S]{
		constructor: constructor,
		once:        sync.Once{},
	}
}

func GetServiceFromRegistry[S any](r *Registry, name string) S {
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

func SetService[S any](name string, constructor ServiceConstructor[S]) {
	SetServiceFromRegistry[S](defaultRegistry, name, constructor)
}

func GetService[S any](name string) S {
	return GetServiceFromRegistry[S](defaultRegistry, name)
}

func (def *serviceDef[S]) init(r *Registry) {
	def.once.Do(func() {
		def.service = def.constructor(r)
	})
}

func InitRegistry(r *Registry) {
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

func Init() {
	InitRegistry(defaultRegistry)
}
