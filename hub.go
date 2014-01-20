package hub

import "sync"

type Hub struct {
	subscribers map[int][]chan interface{}
	sync.RWMutex
}

func New() *Hub {
	return &Hub{subscribers: make(map[int][]chan interface{})}
}

func (h *Hub) Sub(eventType int) chan interface{} {
	c := make(chan interface{})
	h.Lock()
	h.subscribers[eventType] = append(h.subscribers[eventType], c)
	h.Unlock()
	return c
}

func (h *Hub) Pub(eventType int, event interface{}) {
	h.RLock()
	if subscribers, ok := h.subscribers[eventType]; ok {
		for _, c := range subscribers {
			c <- event
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

var defaultHub *Hub

func init() {
	defaultHub = New()
}

func Sub(eventType int) chan interface{} {
	return defaultHub.Sub(eventType)
}

func Pub(eventType int, event interface{}) {
	defaultHub.Pub(eventType, event)
}
