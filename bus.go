package bus

import (
	"context"
	"sync"
	"unsafe"
)

// Creates a new bus with a worker.
func NewBus(ctx context.Context, capacity int) *Bus {
	b := &Bus{
		ctx:      ctx,
		queue:    make(chan *event, capacity),
		subQueue: make(chan sub, 2),
		subs:     make(map[subKey][]unsafe.Pointer),
	}

	b.startWorker()

	return b
}

// Event bus
type Bus struct {
	ctx       context.Context
	queue     chan *event
	subQueue  chan sub
	subs      map[subKey][]unsafe.Pointer
	eventPool eventPool
	wg        sync.WaitGroup
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

func (b *Bus) Close() {
	close(b.queue)
	b.wg.Wait()
}
