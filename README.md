# bus
A generic, thread-safe, highly optimized event bus for Go.

## Install
```
go get github.com/webmafia/bus
```

## Usage: Simple
```go
const (
	Created bus.Topic = iota
	Updated
	Deleted
)

type User struct {
	Id   int
	Name string
}

// Create a new event bus with a buffer of 32 events
b := bus.NewBus(32)

// Subscribe for created users
bus.Sub(b, Created, func(u *User) {
	log.Printf("User %d was created", u.Id)
})

// Subscribe for updated users
bus.Sub(b, Updated, func(u *User) {
	log.Printf("User %d was updated", u.Id)
})

// Publish to Created topic. This will always do 1 allocation due to the pointer.
bus.Pub(b, Created, &User{
	Id:   123,
	Name: "John Doe"
})
```

## Usage: Pool
```go
const (
	Created bus.Topic = iota
	Updated
	Deleted
)

type User struct {
	Id   int
	Name string
}

// Create a new event bus with a buffer of 32 events
b := bus.NewBus(32)

// Subscribe for created users
bus.Sub(b, created, func(u *User) {
	log.Printf("User %d was created", u.Id)
})

// Subscribe for updated users
bus.Sub(b, Updated, func(u *User) {
	log.Printf("User %d was updated", u.Id)
})

userPool := sync.Pool{
	New: func() any {
		return new(User)
	},
}

// Acquire object from pool
user := userPool.Get().(*User)

// Set properties
user.Id = 123
user.Name = "John Doe"

// Publish to Created topic. The object will be returned to the pool after
// handling all subscriptions.
bus.PubPool(b, Created, user, &userPool)
```