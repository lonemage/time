package time

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type xcloudDriver struct {
	addr string

	cb func(string, Time)

	c *client
}

func NewXcloudDriver(addr string, cb func(string, Time)) (Driver, error) {
	con, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	c := newClient(con)
	return &xcloudDriver{addr: addr,
		cb: cb,
		c:  c,
	}, nil
}

type SyncRequest struct {
	Ns       string `json:"ns"`
	Local    int64  `json:"local"`
	Now      int64  `json:"now,omitempty"`
	Location string `json:"location,omitempty"`
}

type SyncResponse struct {
	Local    int64  `json:"local"`
	Now      int64  `json:"now"`
	Location string `json:"location"`
}

func (x *xcloudDriver) Sync(ns string) (Time, error) {
	var req = struct {
		Ns string `json:"ns"`
	}{Ns: ns}

	var resp = struct {
		Now int64  `json:"now"`
		Loc string `json:"loc"`
	}{}

	if err := x.call(&req, &resp); err != nil {
		return Time{}, err
	}

	loc, err := LoadLocation(resp.Loc)
	if err != nil {
		return Time{}, err
	}

	return Unix(resp.Now/int64(time.Second), resp.Now%int64(time.Second)).In(loc), nil
}

func (x *xcloudDriver) In(ns string, loc *Location) error {
	var req = struct {
		Ns  string `json:"ns"`
		Loc string `json:"loc"`
	}{Ns: ns, Loc: loc.String()}

	return x.call(&req, nil)
}

func (x *xcloudDriver) SetTime(ns string, now Time) error {
	var req = struct {
		Ns  string `json:"ns"`
		Now int64  `json:"now"`
	}{Ns: ns, Now: now.UnixNano()}

	return x.call(&req, nil)
}

func (x *xcloudDriver) call(req, resp interface{}) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}

	if x.c == nil {
		con, err := net.Dial("tcp", x.addr)
		if err != nil {
			return err
		}
		c := newClient(con)
		x.c = c
	}

	data, err = x.c.call(data)
	if err != nil {
		return err
	}

	if resp != nil {
		return json.Unmarshal(data, resp)
	}
	return nil
}

type client struct {
	pid    uint16
	con    net.Conn
	endian binary.ByteOrder
	ch     chan []byte
	calls  sync.Map // map[uint16]chan []byte
}

func newClient(con net.Conn) *client {
	c := &client{
		con:    con,
		endian: binary.LittleEndian,
		ch:     make(chan []byte, 16),
		calls:  sync.Map{}, //  make(map[uint16]chan []byte, 16),
	}
	go c.read()
	return c
}

type Packet struct {
	Length uint16
	ID     uint16
	Data   []byte
}

func (c *client) call(b []byte) ([]byte, error) {
	if len(b) == 0 {
		return nil, errors.New("zero length")
	}

	c.pid++
	pid := c.pid

	var p Packet
	p.Length = uint16(len(b))
	p.ID = pid
	p.Data = b
	if _, err := c.con.Write(b); err != nil {
		return nil, err
	}

	ch := make(chan []byte, 1)
	c.calls.Store(pid, ch)
	select {
	case data := <-ch:
		return data, nil
	case <-time.After(time.Second * 30):
		return nil, errors.New("timeout")
	}
}

func (c *client) write(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	var p Packet
	p.Length = uint16(len(b))
	p.Data = b
	_, err := c.con.Write(b)
	return err
}

func (c *client) read() error {
	var buf [1024]byte
	var pos int
	for {
		n, err := c.con.Read(buf[pos:])
		if err != nil {
			return err
		}
		pos += n

		for pos >= 4 {
			length := c.endian.Uint16(buf[:])
			if length > 1000 {
				return fmt.Errorf("invalid length %d", length)
			}
			if pos < int(length)+4 {
				break
			}

			var p Packet
			p.Length = length
			p.ID = c.endian.Uint16(buf[2:])
			p.Data = make([]byte, p.Length)
			copy(p.Data, buf[4:p.Length+4])
			pos -= int(p.Length) + 4

			if p.ID == 0 {
				c.ch <- p.Data
			} else if ch, ok := c.calls.Load(p.ID); !ok {
				fmt.Errorf("call %d disappear", p.ID)
			} else {
				ch.(chan []byte) <- p.Data
			}
		}
	}
}

func (c *client) close() {
	// TODO
	// closeread
}
