package main

import (
	"fmt"
	"time"

	"github.com/tomr-ninja/indi"
)

// Dependency tree:
// A -> B -> C
// A -> D

type (
	A struct{ x int }
	B struct{ x int }
	C struct{ x int }
	D struct{ x int }

	IB interface {
		B() int
	}
	IC interface {
		C() int
	}
	ID interface {
		D() int
	}
)

func (b *B) B() int {
	return b.x
}

func (c *C) C() int {
	return c.x
}

func (d *D) D() int {
	return d.x
}

// NOTE: constructors accept interfaces, not concrete types

func NewA(b IB, d ID) (*A, error) {
	time.Sleep(time.Second)

	return &A{x: b.B() + d.D()}, nil
}

func NewB(c IC) (*B, error) {
	time.Sleep(time.Second)

	return &B{x: c.C() * 2}, nil
}

func NewC() (*C, error) {
	time.Sleep(time.Second)

	return &C{x: 30}, nil
}

func NewD() (*D, error) {
	time.Sleep(time.Second)

	return &D{x: 30}, nil
}

func main() {
	var (
		// NOTE: actual types here
		a A
		b B
		c C
		d D
	)

	indi.Declare(&a, func() (*A, error) { return NewA(&b, &d) }, &b, &d)
	indi.Declare(&b, func() (*B, error) { return NewB(&c) }, &c)
	indi.Declare(&c, NewC)
	indi.Declare(&d, NewD)

	now := time.Now()
	if err := indi.Init(); err != nil {
		panic(err)
	}

	fmt.Println(time.Since(now)) // 3 seconds, not 4
	fmt.Println(a, b, c, d)      // {90} {60} {30} {30}
}
