package ace

import (
	"sync"
	"testing"
	"time"
)

type fakeHost struct {
	mu    sync.Mutex
	calls []string
	out   []byte
}

func (f *fakeHost) add(name string) {
	f.mu.Lock()
	f.calls = append(f.calls, name)
	f.mu.Unlock()
}

func (f *fakeHost) Heartbeat() error          { f.add("heartbeat"); return nil }
func (f *fakeHost) ProcessReceived() error    { f.add("process"); return nil }
func (f *fakeHost) SendStatus() error         { f.add("status"); return nil }
func (f *fakeHost) DetectSpeedHack(int) error { f.add("speed"); return nil }
func (f *fakeHost) PullAnti() ([]byte, error) { f.add("poll"); return f.out, nil }
func (f *fakeHost) PushAnti([]byte) error     { f.add("feed"); return nil }

func (f *fakeHost) has(name string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, c := range f.calls {
		if c == name {
			return true
		}
	}
	return false
}

func (f *fakeHost) count(name string) int {
	f.mu.Lock()
	defer f.mu.Unlock()
	n := 0
	for _, c := range f.calls {
		if c == name {
			n++
		}
	}
	return n
}

func TestStepOfficialIntervals(t *testing.T) {
	h := &fakeHost{out: []byte("ace")}
	var got []byte
	life := New(h, func(data []byte) ([]byte, error) {
		got = append([]byte{}, data...)
		return []byte("srv"), nil
	})
	base := time.Unix(0, 0)
	life.Reset(base)
	life.Step(base.Add(5 * time.Second))
	life.Step(base.Add(25 * time.Second))
	life.Step(base.Add(30 * time.Second))
	life.Step(base.Add(150 * time.Second))
	if !h.has("heartbeat") || !h.has("status") || !h.has("speed") || !h.has("feed") {
		t.Fatalf("calls %v", h.calls)
	}
	if string(got) != "ace" {
		t.Fatalf("upload %q", got)
	}
	if snap := life.Snapshot(); snap.Uploads < 1 || snap.Reports < 1 {
		t.Fatalf("status %+v", snap)
	}
}

func TestStatusRepeatsAndSkipsFuncCheck(t *testing.T) {
	h := &fakeHost{}
	life := New(h, func([]byte) ([]byte, error) { return nil, nil })
	base := time.Unix(0, 0)
	life.Reset(base)
	for _, d := range []time.Duration{150, 180, 300, 450, 600} {
		life.Step(base.Add(d * time.Second))
	}
	if n := h.count("status"); n != 4 {
		t.Fatalf("status %d", n)
	}
	if h.has("functions") {
		t.Fatal("不应该发函数完整性检查")
	}
}
