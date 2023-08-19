package main

import (
	"github.com/tomr-ninja/di"
	"github.com/tomr-ninja/di/example/services"
	"time"
)

// Dependency tree:
// A -> B -> C
// A -> D

func main() {
	// 10 seconds
	indi.SetService("A", func(r *indi.Registry) *services.A {
		return services.NewServiceA(
			indi.GetServiceFromRegistry[*services.B](r, "B"),
			indi.GetServiceFromRegistry[*services.D](r, "D"),
		)
	})
	// 15 seconds
	indi.SetService("B", func(r *indi.Registry) *services.B {
		return services.NewServiceB(indi.GetServiceFromRegistry[*services.C](r, "C"))
	})
	// 5 seconds
	indi.SetService("C", func(r *indi.Registry) *services.C {
		return services.NewServiceC()
	})
	// 10 seconds
	indi.SetService("D", func(r *indi.Registry) *services.D {
		return services.NewServiceD()
	})

	now := time.Now()
	indi.Init()
	println(time.Since(now).String()) // should be 30 seconds, not 40
}
