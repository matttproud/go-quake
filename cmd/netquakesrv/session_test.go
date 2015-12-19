package main

import (
	"bytes"
	"errors"
	"io"
	"net"
	"reflect"
	"testing"
	"time"
)

type fakeConn struct {
	net.PacketConn
	Closes           []func() error
	SetReadDeadlines []func(time.Time) error
	ReadFroms        []func([]byte) (int, net.Addr, error)
}

func (c *fakeConn) Close() error {
	err := c.Closes[0]()
	c.Closes = c.Closes[1:]
	return err
}
func (c *fakeConn) SetReadDeadline(t time.Time) error {
	err := c.SetReadDeadlines[0](t)
	c.SetReadDeadlines = c.SetReadDeadlines[1:]
	return err
}
func (c *fakeConn) ReadFrom(data []byte) (int, net.Addr, error) {
	n, addr, err := c.ReadFroms[0](data)
	c.ReadFroms = c.ReadFroms[1:]
	return n, addr, err
}

type testError struct {
	error
	timeout, temporary bool
}

func (t testError) Timeout() bool   { return t.timeout }
func (t testError) Temporary() bool { return t.temporary }

func TestReadDatagram(t *testing.T) {
	temporaryErr := testError{
		error:     errors.New("test temporary"),
		temporary: true,
	}
	permanentErr := testError{
		error:     errors.New("test permanent"),
		temporary: false,
	}
	timeoutErr := testError{
		error:   errors.New("test timeout"),
		timeout: true,
	}
	for _, test := range []struct {
		f    func([]byte) (int, net.Addr, error)
		read []byte
		err  error
	}{
		{
			f: func(d []byte) (int, net.Addr, error) {
				return 0, nil, nil
			},
			read: []byte(nil),
			err:  errShortRead,
		},
		{
			f: func(d []byte) (int, net.Addr, error) {
				copy(d, []byte{0})
				return 1, nil, nil
			},
			read: []byte(nil),
			err:  errShortRead,
		},
		{
			f: func(d []byte) (int, net.Addr, error) {
				copy(d, []byte{0, 1, 2, 3, 4, 5, 6, 7})
				return 8, nil, nil
			},
			read: []byte{0, 1, 2, 3, 4, 5, 6, 7},
			err:  nil,
		},
		{
			f: func(d []byte) (int, net.Addr, error) {
				copy(d, []byte{0, 1, 2, 3, 4, 5, 6, 7, 8})
				return 9, nil, nil
			},
			read: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8},
			err:  nil,
		},
		{
			f: func(d []byte) (int, net.Addr, error) {
				return 0, nil, temporaryErr
			},
			read: []byte(nil),
			err:  temporaryErr,
		},
		{
			f: func(d []byte) (int, net.Addr, error) {
				return 0, nil, permanentErr
			},
			read: []byte(nil),
			err:  permanentErr,
		},
		{
			f: func(d []byte) (int, net.Addr, error) {
				return 0, nil, timeoutErr
			},
			read: []byte(nil),
			err:  timeoutErr,
		},
	} {
		conn := &fakeConn{
			ReadFroms: []func([]byte) (int, net.Addr, error){test.f},
		}
		out, err := readDatagram(conn, nil)
		if got, want := out, test.read; !bytes.Equal(got, want) {
			t.Errorf("got = %v, want = %v", got, want)
		}
		if got, want := err, test.err; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
	}
}

/*
func TestSession(t *testing.T) {
	var canceled, removed bool
	localAddr, err := net.ResolveIPAddr("ip", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	remoteAddr, err := net.ResolveIPAddr("ip", "8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	conn := fakeConn{
		Closes:           []func() error{func() error { return nil }},
		SetReadDeadlines: []func(time.Time) error{func(time.Time) error { return nil }},
		ReadFroms: []func(time.Time) error{func(time.Time) error { return nil }},
	}
	ses := &Session{
		Cancel:        func() { canceled = true },
		Remove:        func() { removed = true },
		Id:            "test-session",
		LocalAddr:     localAddr,
		LocalPort:     31337,
		RemoteAddr:    remoteAddr,
		Conn:          new(fakeConn),
		disconnectSig: make(chan struct{}),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if got, want := ses.Loop(ctx), context.Canceled; got != want {
		t.Errorf("got = %v, want = %v", got, want)
	}
	if got, want := canceled, true; got != want {
		t.Errorf("got = %v, want = %v", got, want)
	}
	if got, want := removed, true; got != want {
		t.Errorf("got = %v, want = %v", got, want)
	}
}
*/

func TestDecodeDatagram(t *testing.T) {
	for _, test := range []struct {
		data  []byte
		dgram *datagram
		err   error
	}{
		{
			data:  []byte(nil),
			dgram: nil,
			err:   io.EOF,
		},
		{
			data:  []byte{0},
			dgram: nil,
			err:   io.ErrUnexpectedEOF,
		},
		{
			data: []byte{0, 16, 0, 8, 0, 0, 0, 1},
			dgram: &datagram{
				seq:   1,
				flags: 1048576,
				data:  []byte{},
			},
			err: nil,
		},
		{
			data: []byte{0, 16, 0, 9, 0, 0, 0, 7, 2},
			dgram: &datagram{
				seq:   7,
				flags: 1048576,
				data:  []byte{2},
			},
			err: nil,
		},
		{
			data:  []byte{0, 16, 0, 9, 0, 0, 0, 7},
			dgram: nil,
			err:   io.EOF,
		},
	} {
		dgram, err := decodePacketBuf(test.data)
		if got, want := dgram, test.dgram; !reflect.DeepEqual(got, want) {
			t.Errorf("got = %v, want = %v", got, want)
		}
		if got, want := err, test.err; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
	}
}

func TestDatagram(t *testing.T) {
	for _, test := range []struct {
		dgram        datagram
		len          int
		seq          int
		isNetCtrl    bool
		data         []byte
		isUnreliable bool
	}{
		{

			dgram: datagram{
				seq:   1,
				data:  []byte{},
				flags: ^(netflagControl | netflagUnreliable),
			},
			seq:          1,
			len:          0,
			isNetCtrl:    false,
			data:         []byte{},
			isUnreliable: false,
		},
		{

			dgram: datagram{
				seq:   2,
				data:  []byte{1},
				flags: netflagControl | netflagUnreliable,
			},
			seq:          2,
			len:          1,
			isNetCtrl:    true,
			data:         []byte{1},
			isUnreliable: true,
		},
		{

			dgram: datagram{
				seq:   3,
				data:  []byte{1, 2},
				flags: ^netflagControl | netflagUnreliable,
			},
			seq:          3,
			len:          2,
			isNetCtrl:    false,
			data:         []byte{1, 2},
			isUnreliable: true,
		},
		{

			dgram: datagram{
				seq:   4,
				data:  []byte{1, 2, 3},
				flags: netflagControl | ^netflagUnreliable,
			},
			seq:          4,
			len:          3,
			isNetCtrl:    true,
			data:         []byte{1, 2, 3},
			isUnreliable: false,
		},
	} {
		if got, want := test.dgram.Len(), test.len; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
		if got, want := test.dgram.Seq(), test.seq; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
		if got, want := test.dgram.IsNetCtrl(), test.isNetCtrl; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
		if got, want := test.dgram.Data(), test.data; !bytes.Equal(got, want) {
			t.Errorf("got = %v, want = %v", got, want)
		}
		if got, want := test.dgram.IsUnreliable(), test.isUnreliable; got != want {
			t.Errorf("got = %v, want = %v", got, want)
		}
	}
}
