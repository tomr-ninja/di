package indi_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/tomr-ninja/indi"
)

type testService struct {
	ready bool
	deps  []*testService
}

func createTestService(d time.Duration, deps ...*testService) (*testService, error) {
	time.Sleep(d)

	return &testService{true, deps}, nil
}

func TestInit(t *testing.T) {
	// Dependency tree:
	// T1 -> T2 -> T3
	// |
	// +---> T4 -> T5
	//       |
	//       +---> T6

	indi.SetService("T1", func(r *indi.Registry) (*testService, error) {
		return createTestService(
			10*time.Millisecond,
			indi.GetServiceFromRegistry[*testService](r, "T2"),
			indi.GetServiceFromRegistry[*testService](r, "T4"),
		)
	})
	indi.SetService("T2", func(r *indi.Registry) (*testService, error) {
		return createTestService(
			10*time.Millisecond,
			indi.GetServiceFromRegistry[*testService](r, "T3"),
		)
	})
	indi.SetService("T3", func(r *indi.Registry) (*testService, error) {
		return createTestService(
			10 * time.Millisecond,
		)
	})
	indi.SetService("T4", func(r *indi.Registry) (*testService, error) {
		return createTestService(
			10*time.Millisecond,
			indi.GetServiceFromRegistry[*testService](r, "T5"),
			indi.GetServiceFromRegistry[*testService](r, "T6"),
		)
	})
	indi.SetService("T5", func(r *indi.Registry) (*testService, error) {
		return createTestService(
			10 * time.Millisecond,
		)
	})
	indi.SetService("T6", func(r *indi.Registry) (*testService, error) {
		return createTestService(
			10 * time.Millisecond,
		)
	})

	start := time.Now()
	indi.Init()
	spent := time.Since(start)

	if spent > 32*time.Millisecond { // +2ms threshold
		t.Errorf("init was supposed to finish by ~30ms, but actually took %v", spent)
	}
}

func TestInit_FailEarly(t *testing.T) {
	// Dependency tree:
	// T1 -> T2
	// T2 fails, no need to wait for T1 then

	indi.SetService("T1", func(r *indi.Registry) (*testService, error) {
		return createTestService(
			10*time.Millisecond,
			indi.GetServiceFromRegistry[*testService](r, "T2"),
		)
	})
	indi.SetService("T2", func(r *indi.Registry) (*testService, error) {
		time.Sleep(time.Millisecond)

		return nil, fmt.Errorf("T2 failed")
	})

	start := time.Now()
	err := indi.Init()
	if spent := time.Since(start); spent > 2*time.Millisecond { // +1ms threshold
		t.Errorf("init was supposed to finish by ~1ms, but actually took %v", spent)
	}
	if err == nil {
		t.Error("expected error")
	}
}

func TestGetService(t *testing.T) {
	t.Parallel()
	indi.SetService("T1", func(r *indi.Registry) (*testService, error) {
		return createTestService(
			10 * time.Millisecond,
		)
	})

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			t.Parallel()

			start := time.Now()
			s := indi.GetService[*testService]("T1")
			if s.ready != true {
				t.Error()
			}
			spent := time.Since(start)
			if spent < 10*time.Millisecond || spent > 12*time.Millisecond { // +2ms threshold
				t.Errorf("init was supposed to finish by ~10ms, but actually took %v", spent)
			}
		})
	}
}

func TestPanic(t *testing.T) {
	t.Run("get unregistered service", func(t *testing.T) {
		assertPanic(t, func(t *testing.T) {
			_ = indi.GetServiceFromRegistry[*testService](indi.NewRegistry(), "T42")
		})
	})

	t.Run("get service with wrong type", func(t *testing.T) {
		type testService2 struct{}
		r := indi.NewRegistry()

		indi.SetServiceFromRegistry[*testService](r, "T1", func(r *indi.Registry) (*testService, error) {
			return createTestService(0)
		})

		assertPanic(t, func(t *testing.T) {
			_ = indi.GetServiceFromRegistry[*testService2](r, "T1")
		})
	})

}

func assertPanic(t *testing.T, f func(*testing.T)) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	f(t)
}
