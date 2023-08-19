package services

type ServiceD struct{}

func NewServiceD() *ServiceD {
	return &ServiceD{}
}

func (d *ServiceD) Ten() int {
	return 10
}
