package xtime

import (
	"sync"
	"time"
)

type Client interface {
	Off(du time.Duration) time.Duration
}

type localClient struct {
	mtx sync.RWMutex
	du  time.Duration
}

func (s *localClient) Off(du time.Duration) time.Duration {
	s.mtx.Lock()
	s.du = du
	s.mtx.Unlock()
	return du
}

type remoteClient struct {
	opts *Options

	addr string
}

func newRemoteClient(opts *Options, addr string) Client {
	c := &remoteClient{opts: opts, addr: addr}
	go c.init()
	return c
}

func (s *remoteClient) init() {
	for {

	}
}

func (s *remoteClient) Off(du time.Duration) time.Duration {
	return du
}
