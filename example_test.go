package hub

import (
	"fmt"
)

// Different event kinds
const (
	happenedA = iota
	happenedB
	happenedC
)

// Our custom event type
type EventA struct {
	arg1, arg2 int
}

// Implement hub.Event interface
func (e EventA) Kind() int { return happenedA }

func ExampleHub() {
	h := New()

	h.Subscribe(happenedA, func(e Event) {
		a := e.(EventA) // Cast to concrete type
		fmt.Println(a.arg1 + a.arg2)
	})

	h.Publish(EventA{2, 3})
	// Output: 5
}
