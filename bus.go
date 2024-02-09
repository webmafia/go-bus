package bus

import (
	"sync"
	"unsafe"
)

// Creates a new bus with one worker. This should be enought in most cases. If the worker falls behind,
// spawn more workers with `NewBusWithWorkers`.
func NewBus(capacity int) *Bus {
	return NewBusWithWorkers(capacity, 1)
}

// Creates a new bus with N workers. This should be used when there are many or heavy subscribers.
func NewBusWithWorkers(capacity int, workers int) *Bus {
	b := NewBusWithoutWorker(capacity)

	for i := 0; i < workers; i++ {
		b.SpawnWorker()
	}

	return b
}

// Creates a new bus without any worker. Spawn your own workers using the `SpawnWorker` method. If the capacity
// is filled up without any worker, there will be deadlock. This should only be used in rare cases.
func NewBusWithoutWorker(capacity int) *Bus {
	b := &Bus{
		queue: make(chan event, capacity),
		subs:  make(map[subscription][]func(unsafe.Pointer)),
	}

	return b
}

// Event bus
type Bus struct {
	queue chan event
	subs  map[subscription][]func(unsafe.Pointer)
	mu    sync.RWMutex
	wg    sync.WaitGroup
}

func (b *Bus) Size() int {
	return len(b.queue)
}

func (b *Bus) MaxSize() int {
	return cap(b.queue)
}

func (b *Bus) Empty() bool {
	return len(b.queue) == 0
}

func (b *Bus) Full() bool {
	return len(b.queue) == cap(b.queue)
}

func (b *Bus) Workers() int {
	return int(waitgroupCount(&b.wg))
}

func (b *Bus) Close() {
	close(b.queue)
	b.wg.Wait()
}
