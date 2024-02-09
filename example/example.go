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
	log.Println("Workers:", b.Workers())
	b.SpawnWorker()
	log.Println("Workers:", b.Workers())

	bus.Sub(b, 0, func(t *Foobar) {
		// Pretend doing some slow work
		time.Sleep(time.Millisecond * 10)
		log.Println("sub A:", t.Id)
	})

	bus.Sub(b, 0, func(t *Foobar) {
		// Pretend doing some slow work
		time.Sleep(time.Millisecond * 10)
		log.Println("sub B:", t.Id)
	})

	for i := 1; i <= 10; i++ {
		log.Println("pub:", i)
		bus.Pub(b, 0, &Foobar{Id: i})
	}

	b.Close()
	log.Println("Workers:", b.Workers())
}
