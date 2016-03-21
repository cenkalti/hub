package hub

import "testing"

const testKind Kind = 1
const testValue = "foo"

type testEvent string

func (e testEvent) Kind() Kind {
	return testKind
}

func TestPubSub(t *testing.T) {
	var h Hub
	var s string

	h.Subscribe(testKind, func(e Event) { s = string(e.(testEvent)) })
	h.Publish(testEvent(testValue))

	if s != testValue {
		t.Errorf("invalid value: %s", s)
	}
}

func TestCancel(t *testing.T) {
	var h Hub
	var called int

	cancel := h.Subscribe(testKind, func(e Event) { called += 1 })
	h.Publish(testEvent(testValue))
	cancel()
	h.Publish(testEvent(testValue))

	if called != 1 {
		t.Error("unexpected call")
	}
}
