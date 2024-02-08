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
		go b.Worker()
	}

	return b
}

// Creates a new bus without any worker. Spawn your own workers using the `Worker` method. If the capacity
// is filled up without any worker, there will be deadlock. This should only be used in rare cases.
func NewBusWithoutWorker(capacity int) *Bus {
	b := &Bus{
		queue: make(chan event, capacity),
		subs:  make(map[subscription][]func(unsafe.Pointer)),
	}

	return b
}

func Pub[T any](b *Bus, topic Topic, v *T) bool {
	return b._Pub(topic, v)
}

func Sub[T any](b *Bus, topic Topic, cb func(*T)) {
	var v *T
	in := toIface(v)
	b._Sub(topic, in.tab, unsafe.Pointer(&cb))
}

type event struct {
	msg unsafe.Pointer
	sub subscription
}

type subscription struct {
	topic Topic
	typ   unsafe.Pointer
}

type SubFn[T any] func(v *T)

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

func (b *Bus) Close() {
	close(b.queue)
}

func (b *Bus) _Sub(topic Topic, typ unsafe.Pointer, cb unsafe.Pointer) {
	sub := subscription{
		topic: topic,
		typ:   typ,
	}

	b.mu.Lock()
	b.subs[sub] = append(b.subs[sub], *(*func(unsafe.Pointer))(cb))
	b.mu.Unlock()
}

// Spawns a new worker. This function is blocking, and should be run in a goroutine.
func (b *Bus) Worker() {
	for {
		ev, ok := <-b.queue

		if !ok {
			return
		}

		b.mu.RLock()
		for _, sub := range b.subs[ev.sub] {
			sub(ev.msg)
		}
		b.mu.RUnlock()
	}
}

func (b *Bus) _Pub(topic Topic, msg any) bool {
	in := toIface(msg)

	b.queue <- event{
		msg: in.data,
		sub: subscription{
			topic: topic,
			typ:   in.tab,
		},
	}

	return true
}
