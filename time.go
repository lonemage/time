package time

import (
	"sync"
	"time"

	_ "unsafe"

	"github.com/cenkalti/backoff/v4"
)

// runtimeNano returns the current value of the runtime clock in nanoseconds.
//go:linkname runtimeNano runtime.nanotime
func runtimeNano() int64

type xTime struct {
	errHandlers []func(error)

	mtx sync.RWMutex

	d Driver

	now      Time
	syncNano int64
}

var gTime *xTime

func Conf(opt ...Option) error {
	var opts = DefaultOptions()
	for _, o := range opt {
		o(opts)
	}

	var xt xTime
	xt.errHandlers = opts.ErrHandlers
	xt.now = time.Now().In(opts.Location)
	if opts.RemoteServerAddr != "" {
		xt.d = NewXcDriver(opts.RemoteServerAddr, xt.onSync)
	} else {
		xt.d = NewFakeDriver()
	}

	if err := xt.syncOnce(); err != nil {
		if opts.RemoteServerForce {
			return err
		} else {
			bo := backoff.NewExponentialBackOff()
			bo.InitialInterval = time.Second * 30
			go func() {
				if err := backoff.Retry(xt.syncOnce, bo); err != nil {
					xt.onErr(err)
				}
			}()
		}
	}

	gTime = &xt
	return nil
}

func (x *xTime) set(v *dValue) error {
	loc, err := time.LoadLocation(v.location)
	if err != nil {
		return err
	}

	now := time.Unix(v.now/int64(Second), v.now%int64(Second)).In(loc)

	x.mtx.Lock()
	defer x.mtx.Unlock()
	x.syncNano = v.monotonic
	x.now = now
	return nil
}

func (x *xTime) syncOnce() error {
	v, err := x.d.Get()
	if err != nil {
		return err
	}

	return x.set(v)
}

func (x *xTime) onSync(v *dValue) {
	if err := x.set(v); err != nil {
		x.onErr(err)
	}
}

func Now() Time { return gTime.Now() }
func (x *xTime) Now() Time {
	nano := runtimeNano()

	x.mtx.RLock()
	nano -= x.syncNano
	now := x.now
	x.mtx.RUnlock()

	return now.Add(time.Duration(nano) * Nanosecond)
}

func Since(tm Time) Duration { return gTime.Since(tm) }
func (x *xTime) Since(tm Time) Duration {
	return x.Now().Sub(tm)
}

func GetLocation() *time.Location { return gTime.GetLocation() }
func (x *xTime) GetLocation() *time.Location {
	x.mtx.RLock()
	defer x.mtx.RUnlock()
	return x.now.Location()
}

func SetLocation(location string) *time.Location { return gTime.SetLocation(location) }
func (x *xTime) SetLocation(location string) *time.Location {
	l, err := x.d.Set(&dValue{location: location})
	if err != nil {
		x.onErr(err)
		return x.GetLocation()
	}

	if err = x.set(l); err != nil {
		x.onErr(err)
	}

	return x.GetLocation()
}

func Set(t Time) { gTime.Set(t) }
func (x *xTime) Set(t Time) {
	v, err := x.d.Set(&dValue{0, t.UnixNano(), t.Location().String()})
	if err != nil {
		x.onErr(err)
		return
	}

	if err := x.set(v); err != nil {
		x.onErr(err)
	}
}

func Pass(du Duration) { gTime.Pass(du) }
func (x *xTime) Pass(du Duration) {
	x.Set(x.Now().Add(du))
}

func (x *xTime) onErr(e error) {
	for _, cb := range x.errHandlers {
		cb(e)
	}
}

// NewTicker 改变时间不影响ticker
func NewTicker(d Duration) *time.Ticker {
	return time.NewTicker(d)
}

func Parse(layout, value string) (Time, error) {
	return ParseInLocation(layout, value, GetLocation())
}
