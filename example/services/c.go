package services

import (
	"math/rand"
	"time"
)

type C struct{}

func NewServiceC() *C {
	time.Sleep(5 * time.Second)

	return &C{}
}

func (c *C) RandInt() int {
	return rand.Int()
}
