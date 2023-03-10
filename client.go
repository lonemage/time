package time

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type dValue struct {
	monotonic int64
	now       int64
	location  string
}

type Driver interface {
	Start()

	Stop()

	Get() (*dValue, error)

	Set(v *dValue) (*dValue, error)
}

type fakeDriver struct {
	mtx sync.RWMutex

	offset   int64
	location string
}

func NewFakeDriver() Driver {
	return &fakeDriver{location: Local.String()}
}

func (f *fakeDriver) Start() {}

func (f *fakeDriver) Stop() {}

func (f *fakeDriver) Get() (*dValue, error) {
	n := runtimeNano()
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	return &dValue{n, time.Now().UnixNano() + f.offset, f.location}, nil
}

func (f *fakeDriver) Set(v *dValue) (*dValue, error) {
	v.monotonic = runtimeNano()
	f.mtx.Lock()
	f.offset = v.now - time.Now().UnixNano()
	f.location = v.location
	f.mtx.Unlock()
	return v, nil
}

type xcDriver struct {
	addr string

	cb func(*dValue)
}

const (
	// TODO use http get/post
	apiGet = "/get"
	apiSet = "/set"
)

func NewXcDriver(addr string, cb func(*dValue)) Driver {
	c := &xcDriver{addr: addr, cb: cb}
	// TODO
	// sync cb
	return c
}

func (x *xcDriver) Start() {
}

func (x *xcDriver) Stop() {
}

type Request struct {
	Local    int64  `json:"local"`
	Now      int64  `json:"now,omitempty"`
	Location string `json:"location,omitempty"`
}

type Response struct {
	Local    int64  `json:"local"`
	Now      int64  `json:"now"`
	Location string `json:"location"`
}

func (x *xcDriver) Get() (*dValue, error) {
	var req Request
	req.Local = runtimeNano()

	resp, err := x.post(apiGet, &req)
	if err != nil {
		return nil, err
	}

	return &dValue{req.Local, resp.Now - (runtimeNano()-req.Local)/2, resp.Location}, nil
}

type SetMessage struct {
	Now      int64  `json:"now,omitempty"`
	Location string `json:"location,omitempty"`
}

func (x *xcDriver) Set(v *dValue) (*dValue, error) {
	var req Request
	if v.monotonic != 0 {
		req.Local = v.monotonic
	} else {
		req.Local = runtimeNano()
	}
	if v.now != 0 {
		req.Now = v.now
	}
	if v.location != "" {
		req.Location = v.location
	}

	resp, err := x.post(apiSet, &req)
	if err != nil {
		return nil, err
	}

	var now int64
	if resp.Now != 0 {
		now = resp.Now - (runtimeNano()-req.Local)/2
	}
	return &dValue{resp.Local, now, resp.Location}, nil
}

type httpResponse struct {
	Ec   int      `json:"ec"`
	Data Response `json:"data,omitempty"`
}

func (s *xcDriver) post(api string, req *Request) (*Response, error) {
	var bd []byte
	if req != nil {
		var err error
		if bd, err = json.Marshal(req); err != nil {
			return nil, err
		}
	}

	httpResp, err := http.Post(s.addr+api, "", bytes.NewReader(bd))
	if err != nil {
		return nil, err
	}

	defer httpResp.Body.Close()
	bd, err = ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	var r httpResponse
	if err := json.Unmarshal(bd, &r); err != nil {
		return nil, err
	}
	return &r.Data, nil
}
