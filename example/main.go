package main

import (
	"time"

	"github.com/tomr-ninja/indi"
	"github.com/tomr-ninja/indi/example/services"
)

// Dependency tree:
// A -> B -> C
// A -> D

func main() {
	// 10 seconds
	indi.SetService("A", func(r *indi.Registry) (*services.A, error) {
		return services.NewServiceA(
			indi.GetServiceFromRegistry[*services.B](r, "B"),
			indi.GetServiceFromRegistry[*services.D](r, "D"),
		)
	})
	// 15 seconds
	indi.SetService("B", func(r *indi.Registry) (*services.B, error) {
		return services.NewServiceB(indi.GetServiceFromRegistry[*services.C](r, "C"))
	})
	// 5 seconds
	indi.SetService("C", func(r *indi.Registry) (*services.C, error) {
		return services.NewServiceC()
	})
	// 10 seconds
	indi.SetService("D", func(r *indi.Registry) (*services.D, error) {
		return services.NewServiceD()
	})

	now := time.Now()
	if err := indi.Init(); err != nil {
		panic(err)
	}

	println(time.Since(now).Seconds()) // should be 30 seconds, not 40
}
