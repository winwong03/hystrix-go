package hystrix

import (
	"context"
	"fmt"
	"sync"

	"github.com/afex/hystrix-go/hystrix/rolling"
)

type poolMetrics struct {
	Mutex   *sync.RWMutex
	Updates chan poolMetricsUpdate

	Name              string
	MaxActiveRequests *rolling.Number
	Executed          *rolling.Number
}

type poolMetricsUpdate struct {
	activeCount int
}

func newPoolMetrics(ctx context.Context, name string) *poolMetrics {
	m := &poolMetrics{}
	m.Name = name
	m.Updates = make(chan poolMetricsUpdate)
	m.Mutex = &sync.RWMutex{}

	m.Reset()

	go m.Monitor(ctx)

	return m
}

func (m *poolMetrics) Reset() {
	m.Mutex.Lock()
	defer m.Mutex.Unlock()

	m.MaxActiveRequests = rolling.NewNumber()
	m.Executed = rolling.NewNumber()
}

func (m *poolMetrics) Monitor(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("No longer waiting in (m *poolMetrics) Monitor(ctx context.Context)")
			return
		case update, ok := <-m.Updates:
			if !ok {
				return
			}
			m.Mutex.RLock()
			m.Executed.Increment(1)
			m.MaxActiveRequests.UpdateMax(float64(update.activeCount))
			m.Mutex.RUnlock()
		}
	}
}
