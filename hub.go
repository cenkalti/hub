package hub

import "sync"

type Event interface {
	Kind() int
}

// Hub is an event dispatcher, publishes events to the subscribers
// which are subscribed for a specific event type.
type Hub struct {
	subscribers map[int][]chan Event
	sync.RWMutex
}

// New returns pointer to a new Hub.
func New() *Hub {
	return &Hub{subscribers: make(map[int][]chan Event)}
}

// Subscribe for the event of specific kind.
// The caller must receive messages from the retured channel.
// Otherwise, the next Publish() will hang.
func (h *Hub) Subscribe(kind int) chan Event {
	c := make(chan Event)
	h.Lock()
	h.subscribers[kind] = append(h.subscribers[kind], c)
	h.Unlock()
	return c
}

// Publish an event to the subscribers.
func (h *Hub) Publish(e Event) {
	h.RLock()
	if subscribers, ok := h.subscribers[e.Kind()]; ok {
		for _, c := range subscribers {
			c <- e
		}
	}
	h.RUnlock()
}

// Close all channels returned by Subscribe().
// Afther this is called, Publish() will panic.
func (h *Hub) Close() {
	h.Lock()
	for _, subscribers := range h.subscribers {
		for _, ch := range subscribers {
			close(ch)
		}
	}
	h.Unlock()
}

var DefaultHub = New()

func Subscribe(kind int) chan Event {
	return DefaultHub.Subscribe(kind)
}

func Publish(e Event) {
	DefaultHub.Publish(e)
}
