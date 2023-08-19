package services

type (
	ServiceA struct {
		serviceB IB
		serviceD ID
	}
	IB interface {
		EvenRandInt() int
	}
	ID interface {
		Ten() int
	}
)

func NewServiceA(serviceB IB, serviceD ID) *ServiceA {
	return &ServiceA{serviceB: serviceB, serviceD: serviceD}
}

func (a *ServiceA) EvenRandIntUpToTen() int {
	return a.serviceB.EvenRandInt() % a.serviceD.Ten()
}
