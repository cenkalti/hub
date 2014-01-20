package hub

import (
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	h := New()

	c := h.Sub(1)

	go func() {
		h.Pub(1, "foo")
	}()

	select {
	case received := <-c:
		if received != "foo" {
			t.Errorf("invalid value: %s", received)
		}
	case <-time.After(time.Second):
		t.Error("timeout")
	}
}
