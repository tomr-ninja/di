package main

import (
	"fmt"
	"time"

	"github.com/tomr-ninja/indi"
)

// Dependency tree:
// A -> B -> C

type (
	// A depends on IB and lazy-loads it, yet it doesn't even know about B's existence
	// Knowing only about IB, A can't initialize B, as it doesn't know how to do it.
	// So we pass a function that knows how to initialize B to A.
	// All this starts to make sense if you consider A being declared in a different package, where B is not present.
	A struct{ lazyBGetter func() (IB, error) }
	B struct{ x int }
	C struct{ x int }

	RareEventsProcessor interface {
		RarelyCalledMethod() (int, error)
	}
	IB interface {
		B() int
	}
	IC interface {
		C() int
	}
)

func NewA(lazyBGetter func() (IB, error)) (*A, error) {
	return &A{lazyBGetter}, nil
}

func NewB(c IC) (*B, error) {
	time.Sleep(time.Second)

	return &B{x: c.C() * 2}, nil
}

func NewC() (*C, error) {
	time.Sleep(time.Second)

	return &C{x: 10}, nil
}

// RarelyCalledMethod - a method that is rarely called.
// It's not worth to initialize B on every A initialization. So we initialize it lazily.
func (a *A) RarelyCalledMethod() (int, error) {
	b, err := a.lazyBGetter()
	if err != nil {
		return 0, err
	}

	v := b.B()

	return v * 2, nil
}

func (b *B) B() int {
	return b.x
}

func (c *C) C() int {
	return c.x
}

func main() {
	programStart := time.Now()

	// we initialize A here, but not it's dependencies
	var (
		b B
		c C
	)
	indi.Declare(&b, func() (*B, error) { return NewB(&c) }, &c)
	indi.Declare(&c, NewC)
	a, err := NewA(indi.LazyLoad[B, IB](&b))
	if err != nil {
		panic(err)
	}

	// we start a goroutine that will handle rare events
	rareEvents := make(chan struct{})
	rareEventsResults := make(chan int)
	go processRareEvents(a, rareEvents, rareEventsResults)

	// ...
	// way more code here, about the stuff that program usually does
	// ...
	// we didn't spend time on initializing B and C, as we didn't need them
	fmt.Println(time.Since(programStart)) // 0s

	// suddenly, this rare event occurs, and to handle it, we initialize B and C as A now needs them
	rareEvents <- struct{}{}
	result := <-rareEventsResults
	fmt.Println(result)                   // 40
	fmt.Println(time.Since(programStart)) // 2s
}

func processRareEvents(processor RareEventsProcessor, rareEvents <-chan struct{}, results chan<- int) {
	for range rareEvents {
		v, err := processor.RarelyCalledMethod()
		if err != nil {
			panic(err)
		}

		results <- v
	}
}
