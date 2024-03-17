package main

import (
	"context"
	"log"

	"github.com/webmafia/bus"
)

type Foobar struct {
	Id int
}

func main() {
	b := bus.NewBus(context.Background(), 16)

	bus.SubVal(b, 1, func(ctx context.Context, t *Foobar) {
		// Pretend doing some slow work
		// time.Sleep(time.Millisecond * 10)
		log.Println("Sub A:", t.Id)
	})

	bus.SubVal(b, 1, func(ctx context.Context, t *Foobar) {
		// Pretend doing some slow work
		// time.Sleep(time.Millisecond * 10)
		log.Println("Sub B:", t.Id)
	})

	bus.Sub(b, 1, func(ctx context.Context) {
		// Pretend doing some slow work
		// time.Sleep(time.Millisecond * 10)
		log.Println("Sub C without data")
	})

	for i := 1; i <= 10; i++ {
		log.Println("Pub:", i)
		bus.Pub(b, 1)
	}

	b.Close()
}
