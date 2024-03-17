# bus
A generic, thread-safe, highly optimized event bus for Go.

## Install
```
go get github.com/webmafia/bus
```

## Usage
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
b := bus.NewBus(context.Background(), 32)

// Subscribe for created users
bus.SubVal(b, Created, func(ctx context.Context, u *User) {
	log.Printf("User %d was created", u.Id)
})

// Subscribe for updated users
bus.SubVal(b, Updated, func(ctx context.Context, u *User) {
	log.Printf("User %d was updated", u.Id)
})

// Publish to Created topic. This will always do 1 allocation due to the pointer.
bus.SubVal(b, Created, &User{
	Id:   123,
	Name: "John Doe"
})
```
