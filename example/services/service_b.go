package services

import "math"

type (
	ServiceB struct {
		serviceC IC
	}
	IC interface {
		RandInt() int
	}
)

func NewServiceB(serviceC IC) *ServiceB {
	return &ServiceB{serviceC: serviceC}
}

func (b *ServiceB) EvenRandInt() int {
	i := b.serviceC.RandInt()
	if i%2 == 0 {
		return i
	}
	if i == math.MaxInt64 {
		return i - 1
	}

	return i + 1
}
