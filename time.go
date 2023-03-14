package time

import (
	"sync"
	"time"

	_ "unsafe"
	// "github.com/cenkalti/backoff"
)

// runtimeNano returns the current value of the runtime clock in nanoseconds.
//go:linkname runtimeNano runtime.nanotime
func runtimeNano() int64

type TimeSystem struct {
	mtx sync.RWMutex

	d Driver

	ns       string
	now      Time
	syncNano int64
	syncPing int64
}

var DefaultSystem *TimeSystem

func Conf() {
	xt, err := New()
	if err != nil {
		panic(err)
	}
	DefaultSystem = xt
}

func New(opt ...Option) (*TimeSystem, error) {
	var opts = DefaultOptions()
	for _, o := range opt {
		o(opts)
	}

	var xt TimeSystem
	if opts.RemoteServerAddr != "" {
		d, err := NewXcloudDriver(opts.RemoteServerAddr, xt.onSync)
		if err != nil {
			return nil, err
		}
		xt.d = d
	} else {
		xt.d = NewLocalDriver()
	}
	if err := xt.Use(opts.Ns); err != nil {
		return nil, err
	}
	return &xt, nil
}

func Use(ns string) error { return DefaultSystem.Use(ns) }
func (t *TimeSystem) Use(ns string) error {
	monotonic := runtimeNano()
	now, err := t.d.Sync(ns)
	if err != nil {
		return err
	}
	ping := monotonic - runtimeNano()/2
	now = now.Add(time.Duration(ping))

	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.ns = ns
	t.now = now
	t.syncNano = monotonic
	t.syncPing = ping
	return nil
}

func Now() Time { return DefaultSystem.Now() }
func (t *TimeSystem) Now() Time {
	nano := runtimeNano()

	t.mtx.RLock()
	nano -= t.syncNano
	now := t.now
	t.mtx.RUnlock()

	return now.Add(time.Duration(nano) * Nanosecond)
}

func Ns() string { return DefaultSystem.Ns() }
func (t *TimeSystem) Ns() string {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	return t.ns
}

func GetLocation() *time.Location { return DefaultSystem.GetLocation() }
func (t *TimeSystem) GetLocation() *time.Location {
	t.mtx.RLock()
	defer t.mtx.RUnlock()
	return t.now.Location()
}

// Set only set time, location will not change
func Set(now Time) error { return DefaultSystem.Set(now) }

// Set only set time, location will not change
func (t *TimeSystem) Set(now Time) error {
	err := t.d.SetTime(t.Ns(), now)
	if err != nil {
		return err
	}

	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.now = now.In(t.now.Location())
	return nil
}

func In(loc *Location) error { return DefaultSystem.In(loc) }
func (t *TimeSystem) In(loc *Location) error {
	if err := t.d.In(t.Ns(), loc); err != nil {
		return err
	}

	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.now = t.now.In(loc)
	return nil
}

func Pass(d Duration) error { return DefaultSystem.Pass(d) }
func (t *TimeSystem) Pass(d Duration) error {
	return t.Set(t.Now().Add(d))
}

func (t *TimeSystem) onSync(ns string, now Time) {
	if ns == t.Ns() {
		t.Use(ns)
	}
}

type Ticker struct {
	gt *time.Ticker
	C  <-chan Time
	sc chan struct{}
}

func (t *Ticker) Stop() {
	t.sc <- struct{}{}
	t.gt.Stop()
}

func (t *Ticker) Reset(d Duration) {
	t.gt.Reset(d)
}

func Tick(d Duration) <-chan Time  { return NewTicker(d).C }
func NewTicker(d Duration) *Ticker { return DefaultSystem.NewTicker(d) }
func (t *TimeSystem) NewTicker(d Duration) *Ticker {
	gt := time.NewTicker(d)
	c := make(chan Time, 1)
	sc := make(chan struct{}, 1)
	go func() {
		for {
			select {
			case <-gt.C:
				select {
				case c <- t.Now():
				default:
				}
			case <-sc:
				return
			}
		}
	}()
	return &Ticker{gt: gt, C: c, sc: sc}
}
