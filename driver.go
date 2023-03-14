package time

import (
	"sync"
	"time"
)

type Driver interface {
	Sync(ns string) (Time, error)

	In(ns string, loc *Location) error

	SetTime(ns string, now Time) error
}

type offsetLoc struct {
	offset time.Duration
	loc    *Location
}

type localDriver struct {
	mtx sync.RWMutex

	ols map[string]*offsetLoc
}

func NewLocalDriver() Driver {
	return &localDriver{ols: make(map[string]*offsetLoc)}
}

func (l *localDriver) Sync(ns string) (Time, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()

	ol, ok := l.ols[ns]
	if !ok {
		return time.Now(), nil
	}
	return time.Now().Add(ol.offset).In(ol.loc), nil
}

func (l *localDriver) In(ns string, loc *Location) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	ol, ok := l.ols[ns]
	if !ok {
		ol = &offsetLoc{}
		l.ols[ns] = ol
	}

	ol.loc = loc
	return nil
}

func (l *localDriver) SetTime(ns string, now Time) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()

	ol, ok := l.ols[ns]
	if !ok {
		ol = &offsetLoc{loc: Local}
		l.ols[ns] = ol
	}

	ol.offset = time.Now().Sub(now)
	return nil
}
