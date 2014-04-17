package hub

import "testing"

const testKind Kind = 1
const testValue = "foo"

type testEvent string

func (e testEvent) Kind() Kind {
	return testKind
}

func TestPubSub(t *testing.T) {
	var s string

	h := New()
	h.Subscribe(testKind, func(e Event) { s = string(e.(testEvent)) })
	h.Publish(testEvent(testValue))

	if s != testValue {
		t.Errorf("invalid value: %s", s)
	}
}

func BenchmarkPublish(b *testing.B) {
	h := New()
	h.Subscribe(testKind, func(e Event) {})
	for i := 0; i < b.N; i++ {
		h.Publish(testEvent(testValue))
	}
}
