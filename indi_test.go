package indi_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/tomr-ninja/indi"
)

const singleServiceInitTime = time.Millisecond

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

func NewA(b IB, d ID) (*A, error) {
	time.Sleep(singleServiceInitTime)

	return &A{x: b.B() + d.D()}, nil
}

func NewB(c IC) (*B, error) {
	time.Sleep(singleServiceInitTime)

	return &B{x: c.C() * 2}, nil
}

func NewC() (*C, error) {
	time.Sleep(singleServiceInitTime)

	return &C{x: 30}, nil
}

func NewD() (*D, error) {
	time.Sleep(singleServiceInitTime)

	return &D{x: 30}, nil
}

func TestInitDefaultGraph(t *testing.T) {
	var (
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
		t.Fatal(err)
	}
	if time.Since(now) >= 4*singleServiceInitTime {
		t.Fatal("unexpectedly long init time")
	}

	if a.x != 90 || b.x != 60 || c.x != 30 || d.x != 30 {
		t.Fatal("unexpected values")
	}
}

func TestInit(t *testing.T) {
	var (
		a A
		b B
		c C
		d D

		g = make(indi.Graph)
	)

	indi.DeclareOnGraph(g, &a, func() (*A, error) { return NewA(&b, &d) }, &b, &d)
	indi.DeclareOnGraph(g, &b, func() (*B, error) { return NewB(&c) }, &c)
	indi.DeclareOnGraph(g, &c, NewC)
	indi.DeclareOnGraph(g, &d, NewD)

	now := time.Now()
	if err := indi.InitGraph(g); err != nil {
		t.Fatal(err)
	}
	if time.Since(now) >= 4*singleServiceInitTime {
		t.Fatal("unexpectedly long init time")
	}

	if a.x != 90 || b.x != 60 || c.x != 30 || d.x != 30 {
		t.Fatal("unexpected values")
	}
}

func TestFailingConstructor(t *testing.T) {
	var (
		a A
		b B
		c C
		d D

		g = make(indi.Graph)
	)

	indi.DeclareOnGraph(g, &a, func() (*A, error) { return NewA(&b, &d) }, &b, &d)
	indi.DeclareOnGraph(g, &b, func() (*B, error) { return NewB(&c) }, &c)
	indi.DeclareOnGraph(g, &c, NewC)
	indi.DeclareOnGraph(g, &d, func() (*D, error) { return nil, fmt.Errorf("test error") })

	if err := indi.InitGraph(g); err.Error() != "test error" {
		t.Errorf("expected %q, got %q", "test error", err.Error())
	}
}

func TestLoad(t *testing.T) {
	var (
		a A
		b B
		c C
		d D

		g = make(indi.Graph)
	)

	indi.DeclareOnGraph(g, &a, func() (*A, error) { return NewA(&b, &d) }, &b, &d)
	indi.DeclareOnGraph(g, &b, func() (*B, error) { return NewB(&c) }, &c)
	indi.DeclareOnGraph(g, &c, NewC)
	indi.DeclareOnGraph(g, &d, NewD)

	if err := indi.LoadFromGraph(g, &b); err != nil {
		t.Fatal(err)
	}
	if a.x != 0 || b.x != 60 || c.x != 30 || d.x != 0 { // a, d are not initialized
		t.Fatal("unexpected values")
	}
}

func TestLoadDefaultGraph(t *testing.T) {
	var (
		a A
		b B
		c C
		d D
	)

	indi.Declare(&a, func() (*A, error) { return NewA(&b, &d) }, &b, &d)
	indi.Declare(&b, func() (*B, error) { return NewB(&c) }, &c)
	indi.Declare(&c, NewC)
	indi.Declare(&d, NewD)

	if err := indi.Load(&b); err != nil {
		t.Fatal(err)
	}
	if a.x != 0 || b.x != 60 || c.x != 30 || d.x != 0 { // a, d are not initialized
		t.Fatal("unexpected values")
	}
}

func TestLazyLoad(t *testing.T) {
	var (
		a A
		b B
		c C
		d D

		g = make(indi.Graph)
	)

	indi.DeclareOnGraph(g, &a, func() (*A, error) { return NewA(&b, &d) }, &b, &d)
	indi.DeclareOnGraph(g, &b, func() (*B, error) { return NewB(&c) }, &c)
	indi.DeclareOnGraph(g, &c, NewC)
	indi.DeclareOnGraph(g, &d, NewD)

	cb := indi.LazyLoadFromGraph[B, IB](g, &b)
	v, err := cb()
	if err != nil {
		t.Fatal(err)
	}

	if v.B() != 60 {
		t.Fatal("unexpected value")
	}
}

func TestLazyLoadDefaultGraph(t *testing.T) {
	var (
		a A
		b B
		c C
		d D
	)

	indi.Declare(&a, func() (*A, error) { return NewA(&b, &d) }, &b, &d)
	indi.Declare(&b, func() (*B, error) { return NewB(&c) }, &c)
	indi.Declare(&c, NewC)
	indi.Declare(&d, NewD)

	cb := indi.LazyLoad[B, IB](&b)
	v, err := cb()
	if err != nil {
		t.Fatal(err)
	}

	if v.B() != 60 {
		t.Fatal("unexpected value")
	}
}

func TestInvalidConstructor(t *testing.T) {
	var (
		a A
		b B
		c C
		d D

		g = make(indi.Graph)
	)

	indi.DeclareOnGraph(g, &a, func() (*A, error) { return NewA(&b, &d) }, &b, &d)
	indi.DeclareOnGraph(g, &b, func() (*B, error) { return NewB(&c) }, &c)
	indi.DeclareOnGraph(g, &c, NewC)
	indi.DeclareOnGraph(g, &d, func() (*D, error) { return nil, nil }) // must return non-nil value if no error

	if err := indi.InitGraph(g); !errors.Is(err, indi.ErrInvalidConstructor) {
		t.Errorf("expected %v, got %v", indi.ErrInvalidConstructor, err)
	}
}
