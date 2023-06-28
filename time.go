package time

import (
	"errors"
	"sync"
	"time"

	_ "unsafe"
)

func init() {
	if err := Conf(); err != nil {
		panic(err)
	}
}

var (
	ErrUnsupport = errors.New("unsupport")
)

type (
	Time     = time.Time
	Duration = time.Duration
	Location = time.Location
	Month    = time.Month
)

type Ticker interface {
	C() <-chan Time

	Stop()

	Reset(d Duration)
}

type Timer interface {
	C() <-chan Time

	Stop() bool

	Reset(d Duration) bool
}

var (
	// use LocalLocation mostly
	Local = time.Local
	UTC   = time.UTC
)

const (
	Layout      = time.Layout      // "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
	ANSIC       = time.ANSIC       //  "Mon Jan _2 15:04:05 2006"
	UnixDate    = time.UnixDate    //  "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = time.RubyDate    //  "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = time.RFC822      //  "02 Jan 06 15:04 MST"
	RFC822Z     = time.RFC822Z     // "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = time.RFC850      // "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = time.RFC1123     // "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = time.RFC1123Z    //  "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = time.RFC3339     //  "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = time.RFC3339Nano // "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = time.Kitchen     //  "3:04PM"
	// Handy time stamps.
	Stamp      = time.Stamp      // "Jan _2 15:04:05"
	StampMilli = time.StampMilli // "Jan _2 15:04:05.000"
	StampMicro = time.StampMicro //  "Jan _2 15:04:05.000000"
	StampNano  = time.StampNano  // "Jan _2 15:04:05.000000000"

	//
	LayoutYmdHms = "2006-01-02 15:04:05"
)

const (
	Nanosecond  = time.Nanosecond
	Microsecond = time.Microsecond
	Millisecond = time.Millisecond
	Second      = time.Second
	Minute      = time.Minute
	Hour        = time.Hour
)

var (
	LoadLocation    = time.LoadLocation
	Date            = time.Date
	Unix            = time.Unix
	Sleep           = time.Sleep
	ParseInLocation = time.ParseInLocation

	// TODO
	// implement depend on timesystem
	NewTimer  = time.NewTimer
	NewTicker = time.NewTicker
	Tick      = time.Tick
	Afer      = time.After
)

var (
	Parse = func(layout, value string) (Time, error) { return ParseInLocation(layout, value, LocalLocation()) }
	Since = func(tm Time) Duration { return Now().Sub(tm) }
	Pass  = func(d Duration) error { return Set(Now().Add(d)) }
)

type TimeSystem interface {
	Use(ns string) TimeSystem

	Ns() string

	LocalLocation() *time.Location

	Now() Time

	Set(t Time) error

	SetLocalLocation(loc *time.Location) error

	Err() error
}

var defTs TimeSystem

var (
	Use              = func(ns string) TimeSystem { return defTs.Use(ns) }
	Ns               = func() string { return defTs.Ns() }
	LocalLocation    = func() *time.Location { return defTs.LocalLocation() }
	SetLocalLocation = func(loc *time.Location) error { return defTs.SetLocalLocation(loc) }
	Now              = func() Time { return defTs.Now() }
	Set              = func(t Time) error { return defTs.Set(t) }
)

func Conf(opt ...Option) error {
	var opts Options
	for _, o := range opt {
		o(&opts)
	}

	var ts TimeSystem
	if opts.UseSystemTime {
		ts = GoSys()
	} else if opts.RemoteServerAddr == "" {
		ts = LocalSys()
	} else {
		ts = XSys(opts.RemoteServerAddr, opts.RemoteServerToken)
	}
	if err := ts.Err(); err != nil {
		return err
	}

	defTs = ts
	return nil
}

type goTimeSystem struct{}

// GoSys returns a go TimeSystemInterface
func GoSys() TimeSystem {
	return &goTimeSystem{}
}

func (g *goTimeSystem) Use(ns string) TimeSystem {
	return g
}

func (g *goTimeSystem) Ns() string {
	return ""
}

func (g *goTimeSystem) LocalLocation() *time.Location {
	return Local
}

func (g *goTimeSystem) Now() Time {
	return time.Now()
}

func (g *goTimeSystem) Set(t Time) error {
	return ErrUnsupport
}

func (g *goTimeSystem) SetLocalLocation(loc *time.Location) error {
	return ErrUnsupport
}

func (g *goTimeSystem) Err() error {
	return nil
}

type offsetLoc struct {
	offset time.Duration
	loc    *Location
}

type localTimeSystem struct {
	mtx sync.RWMutex

	curNs string

	oss map[string]*offsetLoc
}

func LocalSys() TimeSystem {
	return &localTimeSystem{oss: make(map[string]*offsetLoc)}
}

func (l *localTimeSystem) Use(ns string) TimeSystem {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	l.curNs = ns
	return l
}

func (l *localTimeSystem) Ns() string {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	return l.curNs
}

func (l *localTimeSystem) LocalLocation() *time.Location {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	return l.get().loc
}

func (l *localTimeSystem) Now() Time {
	l.mtx.Lock()
	ol := l.get()
	l.mtx.Unlock()

	return time.Now().In(ol.loc).Add(ol.offset)
}

func (l *localTimeSystem) Set(t Time) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	ol := l.get()
	ol.offset = t.Sub(time.Now())
	ol.loc = t.Location()
	return nil
}

func (l *localTimeSystem) SetLocalLocation(loc *time.Location) error {
	l.mtx.Lock()
	ol := l.get()
	ol.loc = loc
	l.mtx.Unlock()

	return nil
}

func (l *localTimeSystem) Err() error {
	return nil
}

func (l *localTimeSystem) get() *offsetLoc {
	ol, ok := l.oss[l.curNs]
	if !ok {
		ol = &offsetLoc{offset: 0, loc: Local}
		l.oss[l.curNs] = ol
	}
	return ol
}

// runtimeNano returns the current value of the runtime clock in nanoseconds.
//
//go:linkname runtimeNano runtime.nanotime
func runtimeNano() int64

type crossTimeSystem struct {
	mtx sync.RWMutex

	d Driver

	curNs    string
	now      Time
	syncNano int64
	syncPing int64
}

func XSys(addr, token string) TimeSystem {
	var cts crossTimeSystem
	d, err := NewXcloudDriver(addr, cts.onSync)
	if err != nil {
		return &crossTimeSystemUsing{err: err}
	}
	cts.d = d
	return &cts
}

type crossTimeSystemUsing struct {
	*crossTimeSystem
	err error
}

func (c *crossTimeSystemUsing) Err() error {
	return c.err
}

func (c *crossTimeSystem) Use(ns string) TimeSystem {
	monotonic := runtimeNano()
	now, err := c.d.Sync(ns)
	if err != nil {
		return &crossTimeSystemUsing{c, err}
	}
	ping := monotonic - runtimeNano()/2
	now = now.Add(time.Duration(ping))

	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.curNs = ns
	c.now = now
	c.syncNano = monotonic
	c.syncPing = ping
	return nil
}

func (c *crossTimeSystem) Ns() string {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.curNs
}

func (c *crossTimeSystem) LocalLocation() *time.Location {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	return c.now.Location()
}

func (c *crossTimeSystem) Now() Time {
	nano := runtimeNano()

	c.mtx.RLock()
	nano -= c.syncNano
	now := c.now
	c.mtx.RUnlock()

	return now.Add(time.Duration(nano) * Nanosecond)
}

// Set only set time, location will not change
func (c *crossTimeSystem) Set(now Time) error {
	err := c.d.SetTime(c.Ns(), now)
	if err != nil {
		return err
	}

	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.now = now.In(c.now.Location())
	return nil
}

func (c *crossTimeSystem) SetLocalLocation(loc *time.Location) error {
	return c.Set(c.Now().In(loc))
}

func (c *crossTimeSystem) Err() error {
	return nil
}

func (c *crossTimeSystem) onSync(ns string, now Time) {
	if ns == c.Ns() {
		c.Use(ns)
	}
}
