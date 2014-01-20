package hub

import (
	"testing"
	"time"
)

const testKind = 1
const testValue = "foo"

type testEvent string

func (e testEvent) Kind() int {
	return testKind
}

func TestPubSub(t *testing.T) {
	h := New()

	c := h.Subscribe(testKind)

	go func() {
		h.Publish(testEvent(testValue))
	}()

	select {
	case received := <-c:
		if received.(testEvent) != testValue {
			t.Errorf("invalid value: %s", received)
		}
	case <-time.After(time.Second):
		t.Error("timeout")
	}
}
