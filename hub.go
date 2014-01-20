package hub

import "sync"

type Event interface {
	Kind() int
}

type Hub struct {
	subscribers map[int][]chan Event
	sync.RWMutex
}

func New() *Hub {
	return &Hub{subscribers: make(map[int][]chan Event)}
}

func (h *Hub) Sub(kind int) chan Event {
	c := make(chan Event)
	h.Lock()
	h.subscribers[kind] = append(h.subscribers[kind], c)
	h.Unlock()
	return c
}

func (h *Hub) Pub(e Event) {
	h.RLock()
	if subscribers, ok := h.subscribers[e.Kind()]; ok {
		for _, c := range subscribers {
			c <- e
		}
	}
	h.RUnlock()
}

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

func Sub(kind int) chan Event {
	return DefaultHub.Sub(kind)
}

func Pub(e Event) {
	DefaultHub.Pub(e)
}
