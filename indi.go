package indi

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidConstructor = errors.New("constructor must return non-nil value if there is no error")
)

var DefaultGraph = make(Graph)

// DeclareOnGraph - declare a node to initialize on a graph. Non-blocking operation.
func DeclareOnGraph[T any](g Graph, ptr *T, constructor func() (*T, error), deps ...any) {
	cb := func(ptr any) error {
		v, err := constructor()
		if err != nil {
			return err
		} else if v == nil {
			return ErrInvalidConstructor
		}

		*(ptr.(*T)) = *v

		return nil
	}

	g.addNode(ptr, cb, deps...)
}

// LoadFromGraph - initialize a node from a graph and all it's dependencies. Blocking operation.
func LoadFromGraph[T any](g Graph, ptr *T) error {
	return g.ensureReadyNode(formatPointerAddress(ptr))
}

// LazyLoadFromGraph - a helper function, wrapping LoadFromGraph.
// The idea is to initialize something and all it's dependencies without knowing what's the concrete type of it.
func LazyLoadFromGraph[T, I any](g Graph, ptr *T) (cb func() (I, error)) {
	if _, ok := interface{}(ptr).(I); !ok {
		panic(fmt.Sprintf("pointer of type %T cannot be casted to %T", ptr, new(I)))
	}

	return func() (I, error) {
		err := LoadFromGraph(g, ptr)
		cast := interface{}(ptr).(I)

		return cast, err
	}
}

// InitGraph - initialize a graph. It is a blocking operation.
func InitGraph(g Graph) error {
	err := g.init()
	g = make(Graph)

	return err
}

// Declare - declare a node to initialize. Non-blocking operation.
func Declare[T any](ptr *T, constructor func() (*T, error), deps ...any) {
	DeclareOnGraph(DefaultGraph, ptr, constructor, deps...)
}

// Load - initialize a node and all it's dependencies. Blocking operation.
func Load[T any](ptr *T) error {
	return LoadFromGraph(DefaultGraph, ptr)
}

// LazyLoad - a helper function, wrapping Load.
// The idea is to initialize something and all it's dependencies without knowing what's the concrete type of it.
func LazyLoad[T, I any](ptr *T) (cb func() (I, error)) {
	return LazyLoadFromGraph[T, I](DefaultGraph, ptr)
}

// Init - initialize a default graph. Blocking operation.
func Init() error {
	return InitGraph(DefaultGraph)
}
