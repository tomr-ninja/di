package indi

import (
	"errors"
)

var (
	ErrInvalidConstructor = errors.New("constructor must return non-nil value if there is no error")
)

var DefaultGraph = make(Graph)

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

func InitGraph(g Graph) error {
	err := g.init()
	g = make(Graph)

	return err
}

func Declare[T any](ptr *T, constructor func() (*T, error), deps ...any) {
	DeclareOnGraph(DefaultGraph, ptr, constructor, deps...)
}

func Init() error {
	return InitGraph(DefaultGraph)
}
