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
	sr := indi.NewRegistry()

	// 10 seconds
	indi.SetService(sr, "A", func(r *indi.Registry) *services.A {
		return services.NewServiceA(
			indi.GetService[*services.B](r, "B"),
			indi.GetService[*services.D](r, "D"),
		)
	})
	// 15 seconds
	indi.SetService(sr, "B", func(r *indi.Registry) *services.B {
		return services.NewServiceB(indi.GetService[*services.C](r, "C"))
	})
	// 5 seconds
	indi.SetService(sr, "C", func(r *indi.Registry) *services.C {
		return services.NewServiceC()
	})
	// 10 seconds
	indi.SetService(sr, "D", func(r *indi.Registry) *services.D {
		return services.NewServiceD()
	})

	now := time.Now()
	indi.InitAll(sr)
	println(time.Since(now).String()) // should be 30 seconds, not 40
}
