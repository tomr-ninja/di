package services

import "time"

type D struct{}

func NewServiceD() *D {
	time.Sleep(10 * time.Second)

	return &D{}
}

func (d *D) Ten() int {
	return 10
}
