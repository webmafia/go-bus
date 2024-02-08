package main

import (
	"log"
	"time"

	"github.com/webmafia/go-bus"
)

type Foobar struct {
	Id int
}

func main() {
	b := bus.NewBus(16)
	defer b.Close(true)

	bus.Sub(b, 0, func(t *Foobar) {
		// Pretend doing some slow work
		time.Sleep(time.Millisecond * 10)
		log.Println("sub:", t.Id)
	})

	for i := 1; i <= 100; i++ {
		log.Println("pub:", i)
		bus.Pub(b, 0, &Foobar{Id: i})
	}
}
