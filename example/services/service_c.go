package services

import "math/rand"

type ServiceC struct{}

func NewServiceC() *ServiceC {
	return &ServiceC{}
}

func (c *ServiceC) RandInt() int {
	return rand.Int()
}
