package indi

import (
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

type (
	// Graph - dependency graph. Has no public methods as it shouldn't be used directly.
	// However, it is exported to allow custom dependency graphs to be used with
	// DeclareOnGraph and InitGraph methods:
	//
	//	var g = make(indi.Graph)
	//	indi.DeclareOnGraph(g, ...)
	//	...
	//	indi.InitGraph(g)
	Graph     map[string]*graphNode
	graphNode struct {
		mux   sync.Mutex
		ptr   any
		ready bool
		init  func(any) error
		deps  []string
	}
)

func (g Graph) addNode(ptr any, init func(any) error, deps ...any) {
	name := formatPointerAddress(ptr)
	depsNames := make([]string, len(deps))
	for i, dep := range deps {
		depsNames[i] = formatPointerAddress(dep)
	}

	g[name] = &graphNode{
		ptr:   ptr,
		ready: false,
		init:  init,
		deps:  depsNames,
	}
}

func (g Graph) init() error {
	eg := errgroup.Group{}
	for name := range g {
		name := name
		eg.Go(func() error {
			return g.ensureReadyNode(name)
		})
	}

	return eg.Wait()
}

func (g Graph) ensureReadyNode(name string) error {
	node, ok := g[name]
	if !ok {
		return fmt.Errorf("node %s not found", name)
	}

	node.mux.Lock()
	defer node.mux.Unlock()

	if node.ready {
		return nil
	}

	for _, dep := range node.deps {
		if err := g.ensureReadyNode(dep); err != nil {
			return err
		}
	}

	if err := node.init(node.ptr); err != nil {
		return err
	}

	node.ready = true

	return nil
}

func formatPointerAddress(ptr any) string {
	return fmt.Sprintf("%p", ptr)
}
