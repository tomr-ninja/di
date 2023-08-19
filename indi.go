package indi

import (
	"context"
	"errors"
	"sync"

	"golang.org/x/sync/errgroup"
)

var (
	defaultRegistry = NewRegistry()

	errUnregisteredService = errors.New("tried to get unregistered service")
	errWrongDefType        = errors.New("wrong def type")
)

type ServiceConstructor[S any] func(*Registry) (S, error)

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
		panic(errUnregisteredService)
	}

	def, ok := c.(*serviceDef[S])
	if !ok {
		panic(errWrongDefType)
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

func (def *serviceDef[S]) init(r *Registry) (err error) {
	def.once.Do(func() {
		def.service, err = def.constructor(r)
	})

	return err
}

func InitRegistry(r *Registry) error {
	type initable interface {
		init(r *Registry) error
	}

	eg, ctx := errgroup.WithContext(context.Background())
	for _, c := range r.services {
		def := c.(initable)
		eg.Go(func() error {
			ch := make(chan error)
			go func() {
				ch <- def.init(r)
			}()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-ch:
				return err
			}
		})
	}

	return eg.Wait()
}

func Init() error {
	return InitRegistry(defaultRegistry)
}
