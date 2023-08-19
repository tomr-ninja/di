package services

import "time"

type (
	A struct {
		b IB
		d ID
	}
	IB interface {
		EvenRandInt() int
	}
	ID interface {
		Ten() int
	}
)

func NewServiceA(b IB, d ID) (*A, error) {
	time.Sleep(10 * time.Second)

	return &A{b: b, d: d}, nil
}

func (a *A) EvenRandIntUpToTen() int {
	return a.b.EvenRandInt() % a.d.Ten()
}
