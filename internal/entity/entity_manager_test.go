package entity

import (
	"errors"
	"testing"
)

// fakeEntity 用于测试的实体实现
type fakeEntity struct {
	addr  string
	value string
}

func (e *fakeEntity) Address() string { return e.addr }
func (e *fakeEntity) Name() string    { return "fake" }
func (e *fakeEntity) Clone() Entity   { c := *e; return &c }

func (e *fakeEntity) Marshal() ([]byte, error) { return []byte("full:" + e.value), nil }

func (e *fakeEntity) Diff(old Entity) []byte {
	o := old.(*fakeEntity)
	if o.value == e.value {
		return nil
	}

	return []byte("diff:" + e.value)
}

type fakeNotifier struct {
	payloads []string
	methods  []string
	err      error
}

func (n *fakeNotifier) SendNotify(payload []byte, address, method string) error {
	if n.err != nil {
		return n.err
	}

	n.payloads = append(n.payloads, string(payload))
	n.methods = append(n.methods, method)
	return nil
}

func TestCalEntityBodyIncremental(t *testing.T) {
	m := NewEntityManager(true)
	e := &fakeEntity{addr: "game/1", value: "a"}

	// 首次是全量
	payload, addr, err := m.CalEntityBody(e)
	if err != nil {
		t.Fatalf("CalEntityBody failed: %v", err)
	}

	if string(payload) != "full:a" || addr != "game/1" {
		t.Fatalf("first payload = %q, addr = %q", payload, addr)
	}

	// 无变化时不再下发
	payload, _, _ = m.CalEntityBody(e)
	if len(payload) != 0 {
		t.Fatalf("expect empty payload when unchanged, got %q", payload)
	}

	// 有变化时走增量
	e.value = "b"
	payload, _, _ = m.CalEntityBody(e)
	if string(payload) != "diff:b" {
		t.Fatalf("diff payload = %q", payload)
	}
}

func TestSendEntityMethodSelection(t *testing.T) {
	notifier := &fakeNotifier{}

	// 增量模式：$update$
	m := NewEntityManager(true)
	if _, method := m.SendEntity(&fakeEntity{addr: "a", value: "x"}, notifier); method != "$update$" {
		t.Fatalf("method = %q, want $update$", method)
	}

	// 非增量模式：$update_clean$
	m2 := NewEntityManager(false)
	if _, method := m2.SendEntity(&fakeEntity{addr: "a", value: "x"}, notifier); method != "$update_clean$" {
		t.Fatalf("method = %q, want $update_clean$", method)
	}

	// fullPayload 强制 clean
	m3 := NewEntityManager(true)
	if _, method := m3.SendEntity(&fakeEntity{addr: "a", value: "x"}, notifier, true); method != "$update_clean$" {
		t.Fatalf("method = %q, want $update_clean$", method)
	}

	if len(notifier.payloads) != 3 {
		t.Fatalf("notifier got %d payloads", len(notifier.payloads))
	}
}

func TestSendEntityNotifierError(t *testing.T) {
	notifier := &fakeNotifier{err: errors.New("closed")}
	m := NewEntityManager(true)

	payload, method := m.SendEntity(&fakeEntity{addr: "a", value: "x"}, notifier)
	if method != "" {
		t.Fatalf("method = %q, want empty on notify failure", method)
	}

	if len(payload) == 0 {
		t.Fatal("payload should still be returned")
	}
}

func TestResetAndCount(t *testing.T) {
	m := NewEntityManager(true)
	_, _, _ = m.CalEntityBody(&fakeEntity{addr: "a", value: "1"})
	_, _, _ = m.CalEntityBody(&fakeEntity{addr: "b", value: "2"})

	if m.EntityCount() != 2 {
		t.Fatalf("count = %d, want 2", m.EntityCount())
	}

	m.Reset()
	if m.EntityCount() != 0 {
		t.Fatalf("count after reset = %d", m.EntityCount())
	}
}
