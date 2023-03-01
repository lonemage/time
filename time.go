package xtime

import (
	"sync"
	"time"
)

type xTime struct {
	mtx sync.RWMutex

	opts *Options

	c Client

	off time.Duration

	loc *time.Location
}

var xt = &xTime{
	c:   &localClient{},
	loc: Local,
}

func Conf(opt ...Option) {
	var opts Options
	for _, o := range opt {
		o(&opts)
	}

	xt.opts = &opts
	xt.c = newRemoteClient(xt.opts, RemoteServerAddr)
}

func (x *xTime) getOff() time.Duration {
	x.mtx.RLock()
	defer x.mtx.RUnlock()
	return x.off
}

func Now() time.Time { return xt.Now() }
func (x *xTime) Now() time.Time {
	n := time.Now()
	if off := x.getOff(); off != 0 {
		n = n.Add(off)
	}
	return n
}

func Since(tm time.Time) time.Duration { return xt.Since(tm) }
func (x *xTime) Since(tm time.Time) time.Duration {
	return x.Now().Sub(tm)
}

func Location() *time.Location { return xt.Location() }
func (x *xTime) Location() *time.Location {
	x.mtx.RLock()
	defer x.mtx.RUnlock()
	return x.loc
}

func Pass(du time.Duration) { xt.Pass(du) }
func (x *xTime) Pass(du time.Duration) {
	off := x.c.Off(du)

	x.mtx.Lock()
	x.off = off
	x.mtx.Unlock()
}

func Set(tm time.Time) { xt.Set(tm) }
func (x *xTime) Set(tm time.Time) {
	x.Pass(tm.Sub(x.Now()))
}

func Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) time.Time {
	return xt.Date(year, month, day, hour, min, sec, nsec, loc)
}
func (x *xTime) Date(year int, month time.Month, day, hour, min, sec, nsec int, loc *time.Location) time.Time {
	tm := time.Date(year, month, day, hour, min, sec, nsec, loc)
	return tm.In(x.Location()).Add(x.getOff())
}

func NewTicker(d time.Duration) *time.Ticker {
	panic("unimplement")
}

func Parse(layout, value string) (time.Time, error) {
	return ParseInLocation(layout, value, Location())
}

func ParseInLocation(layout, value string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(layout, value, loc)
}

func Unix(sec int64, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}
