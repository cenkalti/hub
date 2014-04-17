// Package hub provides a simple event dispatcher for publish/subscribe pattern.
package hub

import (
	"container/list"
	"sync"
)

type Kind int

// Event is an interface for published events.
type Event interface {
	Kind() Kind
}

// Hub is an event dispatcher, publishes events to the subscribers
// which are subscribed for a specific event type.
type Hub struct {
	subscribers map[Kind]*list.List
	m           sync.RWMutex
}

// New returns pointer to a new Hub.
func New() *Hub {
	return &Hub{
		subscribers: make(map[Kind]*list.List),
	}
}

// Subscribe registers the handler for the event of a specific kind.
func (h *Hub) Subscribe(kind Kind, handler func(Event)) (cancel func()) {
	h.m.Lock()
	defer h.m.Unlock()

	if h.subscribers[kind] == nil {
		h.subscribers[kind] = list.New()
	}
	el := h.subscribers[kind].PushBack(handler)

	return func() {
		h.m.Lock()
		defer h.m.Unlock()

		h.subscribers[kind].Remove(el)
		if h.subscribers[kind].Len() == 0 {
			delete(h.subscribers, kind)
		}
	}
}

// Publish an event to the subscribers.
func (h *Hub) Publish(e Event) {
	h.m.RLock()
	defer h.m.RUnlock()

	if handlers, ok := h.subscribers[e.Kind()]; ok {
		for el := handlers.Front(); el != nil; el = el.Next() {
			el.Value.(func(Event))(e)
		}
	}
}

// DefaultHub is the default Hub used by Publish and Subscribe.
var DefaultHub = New()

// Subscribe registers the handler for the event of a specific kind
// in the DefaultHub.
func Subscribe(kind Kind, handler func(Event)) (cancel func()) {
	return DefaultHub.Subscribe(kind, handler)
}

// Publish an event to the subscribers in DefaultHub.
func Publish(e Event) {
	DefaultHub.Publish(e)
}
