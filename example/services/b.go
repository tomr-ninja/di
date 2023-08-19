package services

import (
	"math"
	"time"
)

type (
	B struct {
		c IC
	}
	IC interface {
		RandInt() int
	}
)

func NewServiceB(c IC) *B {
	time.Sleep(15 * time.Second)

	return &B{c: c}
}

func (b *B) EvenRandInt() int {
	i := b.c.RandInt()
	if i%2 == 0 {
		return i
	}
	if i == math.MaxInt64 {
		return i - 1
	}

	return i + 1
}
