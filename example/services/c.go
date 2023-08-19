package services

import (
	"math/rand"
	"time"
)

type C struct{}

func NewServiceC() (*C, error) {
	time.Sleep(5 * time.Second)

	return &C{}, nil
}

func (c *C) RandInt() int {
	return rand.Int()
}
