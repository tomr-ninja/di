package main

import (
	"runtime"
	"sync"

	"github.com/tomr-ninja/di"
	"github.com/tomr-ninja/di/example/services"
)

// Dependency tree:
// A -> B -> C
// A -> D

func main() {
	sr := di.NewRegistry()

	di.SetService(sr, "serviceA", func(r *di.Registry) *services.ServiceA {
		return services.NewServiceA(
			di.GetService[*services.ServiceB](r, "serviceB"),
			di.GetService[*services.ServiceD](r, "serviceD"),
		)
	})
	di.SetService(sr, "serviceB", func(r *di.Registry) *services.ServiceB {
		return services.NewServiceB(di.GetService[*services.ServiceC](r, "serviceC"))
	})
	di.SetService(sr, "serviceC", func(r *di.Registry) *services.ServiceC {
		return services.NewServiceC()
	})
	di.SetService(sr, "serviceD", func(r *di.Registry) *services.ServiceD {
		return services.NewServiceD()
	})

	serviceA := di.GetService[*services.ServiceA](sr, "serviceA")

	numCPU := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(numCPU)
	for i := 0; i < numCPU; i++ {
		go func() {
			println(serviceA.EvenRandIntUpToTen())
			wg.Done()
		}()
	}

	wg.Wait()
}
