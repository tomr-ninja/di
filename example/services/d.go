package services

import "time"

type D struct{}

func NewServiceD() (*D, error) {
	time.Sleep(10 * time.Second)

	return &D{}, nil
}

func (d *D) Ten() int {
	return 10
}
