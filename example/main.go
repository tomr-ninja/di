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
	indi.Set("A", func(r *indi.Registry) (a *services.A, err error) {
		var (
			b *services.B
			d *services.D
		)
		if b, err = indi.GetFromRegistry[*services.B](r, "B"); err != nil {
			return nil, err
		}
		if d, err = indi.GetFromRegistry[*services.D](r, "D"); err != nil {
			return nil, err
		}

		return services.NewServiceA(b, d) // 10 seconds
	})
	indi.Set("B", func(r *indi.Registry) (b *services.B, err error) {
		var c *services.C
		if c, err = indi.GetFromRegistry[*services.C](r, "C"); err != nil {
			return nil, err
		}

		return services.NewServiceB(c) // 15 seconds
	})
	indi.Set("C", func(r *indi.Registry) (*services.C, error) {
		return services.NewServiceC() // 5 seconds
	})
	indi.Set("D", func(r *indi.Registry) (*services.D, error) {
		return services.NewServiceD() // 10 seconds
	})

	now := time.Now()
	if err := indi.Init(); err != nil {
		panic(err)
	}

	println(time.Since(now).Seconds()) // should be 30 seconds, not 40
}
