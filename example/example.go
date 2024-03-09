package main

import (
	"log"
	"time"

	"github.com/webmafia/bus"
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
		log.Println("Sub A:", t.Id)
	})

	bus.Sub(b, 0, func(t *Foobar) {
		// Pretend doing some slow work
		time.Sleep(time.Millisecond * 10)
		log.Println("Sub B:", t.Id)
	})

	bus.Sub(b, 1, func(t *Foobar) {
		// Pretend doing some slow work
		time.Sleep(time.Millisecond * 10)
		log.Println("Sub C:", t.Id)
	})

	for i := 1; i <= 10; i++ {
		log.Println("Pub:", i)
		bus.Pub(b, bus.Topic(i%2), &Foobar{Id: i})
	}

	b.Close()
	log.Println("Workers:", b.Workers())
}
