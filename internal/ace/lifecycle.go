package ace

import (
	"sync"
	"time"

	"github.com/alttab8520/qqfarm-sdk/internal/crypto"
)

const (
	ProcessInterval   = 5 * time.Second
	HeartbeatInterval = 25 * time.Second
	PollInterval      = 5 * time.Second
	SpeedInterval     = 30 * time.Second
	StatusInterval    = 150 * time.Second
	MaxBackoff        = 30 * time.Second
)

type Uploader func(data []byte) ([]byte, error)

type Status struct {
	Uploads   int    `json:"uploads"`
	Reports   int    `json:"status_reports"`
	Failures  int    `json:"failures"`
	LastError string `json:"last_error,omitempty"`
}

type Life struct {
	host   crypto.Engine
	upload Uploader
	clock  func() time.Time

	mu           sync.Mutex
	stop         chan struct{}
	running      bool
	inFlight     bool
	nextProcess  time.Time
	nextHeart    time.Time
	nextPoll     time.Time
	nextSpeed    time.Time
	nextStatus   time.Time
	lastSpeed    time.Time
	uploads      int
	reports      int
	failures     int
	lastError    string
}

func New(host crypto.Engine, upload Uploader) *Life {
	if host == nil {
		host = crypto.AsEngine(nil)
	}
	return &Life{
		host:   host,
		upload: upload,
		clock:  time.Now,
		stop:   make(chan struct{}),
	}
}

func (l *Life) Reset(now time.Time) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.nextProcess = now.Add(ProcessInterval)
	l.nextHeart = now.Add(HeartbeatInterval)
	l.nextPoll = now.Add(PollInterval)
	l.nextSpeed = now.Add(SpeedInterval)
	l.nextStatus = now.Add(StatusInterval)
	l.lastSpeed = now
	l.failures = 0
	l.lastError = ""
}

func (l *Life) Snapshot() Status {
	l.mu.Lock()
	defer l.mu.Unlock()
	return Status{Uploads: l.uploads, Reports: l.reports, Failures: l.failures, LastError: l.lastError}
}

func (l *Life) Start() {
	l.mu.Lock()
	if l.running {
		l.mu.Unlock()
		return
	}
	l.running = true
	l.stop = make(chan struct{})
	l.mu.Unlock()
	l.Reset(l.clock())
	go l.loop()
}

func (l *Life) Stop() {
	l.mu.Lock()
	if !l.running {
		l.mu.Unlock()
		return
	}
	l.running = false
	close(l.stop)
	l.mu.Unlock()
}

func (l *Life) loop() {
	tick := time.NewTicker(200 * time.Millisecond)
	defer tick.Stop()
	for {
		select {
		case <-l.stop:
			return
		case <-tick.C:
			l.Step(l.clock())
		}
	}
}

func (l *Life) Step(now time.Time) {
	if due(&l.nextProcess, now, ProcessInterval) {
		_ = l.host.ProcessReceived()
	}
	if due(&l.nextHeart, now, HeartbeatInterval) {
		if err := l.host.Heartbeat(); err != nil {
			l.fail(err)
		}
	}
	if due(&l.nextSpeed, now, SpeedInterval) {
		elapsed := int(now.Sub(l.lastSpeed) / time.Millisecond)
		if elapsed < 0 {
			elapsed = 0
		}
		l.lastSpeed = now
		_ = l.host.DetectSpeedHack(elapsed)
	}
	if due(&l.nextStatus, now, StatusInterval) {
		if err := l.host.SendStatus(); err != nil {
			l.fail(err)
		} else {
			l.mu.Lock()
			l.reports++
			l.mu.Unlock()
		}
	}
	l.poll(now)
}

func (l *Life) poll(now time.Time) {
	l.mu.Lock()
	if l.inFlight || now.Before(l.nextPoll) {
		l.mu.Unlock()
		return
	}
	l.inFlight = true
	l.mu.Unlock()
	defer func() {
		l.mu.Lock()
		l.inFlight = false
		l.mu.Unlock()
	}()

	data, err := l.host.PullAnti()
	if err != nil {
		l.backoff(now, err)
		return
	}
	if len(data) == 0 || l.upload == nil {
		l.mu.Lock()
		l.nextPoll = now.Add(PollInterval)
		l.mu.Unlock()
		return
	}
	reply, err := l.upload(data)
	if err != nil {
		l.backoff(now, err)
		return
	}
	if len(reply) > 0 {
		_ = l.host.PushAnti(reply)
	}
	l.mu.Lock()
	l.failures = 0
	l.lastError = ""
	l.uploads++
	l.nextPoll = now.Add(PollInterval)
	l.mu.Unlock()
}

func (l *Life) fail(err error) {
	if err == nil {
		return
	}
	l.mu.Lock()
	l.lastError = err.Error()
	l.mu.Unlock()
}

func (l *Life) backoff(now time.Time, err error) {
	l.mu.Lock()
	l.failures++
	l.lastError = err.Error()
	shift := l.failures
	if shift > 5 {
		shift = 5
	}
	wait := time.Second << uint(shift)
	if wait > MaxBackoff {
		wait = MaxBackoff
	}
	l.nextPoll = now.Add(wait)
	l.mu.Unlock()
}

func due(next *time.Time, now time.Time, interval time.Duration) bool {
	if now.Before(*next) {
		return false
	}
	*next = now.Add(interval)
	return true
}
