package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type SessionBuf struct {
	ops []instruction
	mtx sync.Mutex
}

func SessionId(addr net.Addr) string { return addr.String() }

type errDuplSession string

func (e errDuplSession) Error() string { return "duplicate session: " + string(e) }

type SessionRegistry map[string]*Session

func (r SessionRegistry) Len() int { return len(r) }

func (r SessionRegistry) NewSession(ctx context.Context, addr net.Addr) (*Session, error) {
	id := SessionId(addr)
	if _, ok := r[id]; ok {
		return nil, errDuplSession(id)
	}
	laddr, err := LocalAddr()
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	ses := &Session{
		Cancel:     cancel,
		Id:         id,
		Conn:       conn,
		RemoteAddr: addr,
		LocalAddr:  conn.LocalAddr(),
		LocalPort:  conn.LocalAddr().(*net.UDPAddr).Port,
		Remove: func() {
			log.Printf("Disconnecting %v...", addr)
			r.Disconnect(addr)
			log.Printf("Disconnected %v", addr)
		},
		disconnectSig: make(chan struct{}),
	}
	r[id] = ses
	go func() {
		if err := ses.Loop(ctx); err != nil {
			log.Printf("Closing session from error %v", err)
			return
		}
		log.Printf("Closing session normally")
	}()
	return ses, nil
}

type errUnknownSession string

func (e errUnknownSession) Error() string { return "unknown session: " + string(e) }

func (r SessionRegistry) Disconnect(addr net.Addr) error {
	id := SessionId(addr)
	sess, ok := r[id]
	if !ok {
		return errUnknownSession(id)
	}
	sess.cleanup()
	delete(r, id)
	return nil
}

func (r SessionRegistry) Find(addr net.Addr) (*Session, bool) {
	id := SessionId(addr)
	sess, ok := r[id]
	return sess, ok
}

func (r SessionRegistry) Close() {
	var wg sync.WaitGroup
	for _, s := range r {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.cleanup()
		}()
	}
	wg.Wait()
	for id := range r {
		delete(r, id)
	}
}

func (s *SessionBuf) Add(op instruction) {
	s.mtx.Lock()
	s.ops = append(s.ops, op)
	s.mtx.Unlock()
}

func (s *SessionBuf) Drain(to []instruction) []instruction {
	s.mtx.Lock()
	to = append(to, s.ops...)
	s.ops = s.ops[:0]
	s.mtx.Unlock()
	return to
}

type Session struct {
	Cancel        func()
	Remove        func()
	Id            string
	LocalAddr     net.Addr
	LocalPort     int
	RemoteAddr    net.Addr
	Conn          net.PacketConn
	Seq           int
	Buf           SessionBuf
	cleanupOnce   sync.Once
	disconnectSig chan struct{}
}

func (s *Session) Loop(ctx context.Context) error {
	defer s.cleanup()
	go func() {
		if err := s.loopNet(ctx); err != nil {
			log.Println(err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
			// XXX: should be a part of frame running.
			if err := s.loop(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Session) loop() error {
	for _, o := range s.Buf.Drain([]instruction(nil)) {
		if err := o(); err != nil {
			return err
		}
	}
	return nil
}

func (s *Session) cleanup() {
	s.cleanupOnce.Do(func() {
		s.Cancel()
		if err := s.Conn.Close(); err != nil {
			log.Printf("could not close %s: %v", s, err)
		}
		<-s.disconnectSig
	})
}

type instruction func() error

type errInvalidInstruction string

func (e errInvalidInstruction) Error() string { return "invalid instruction: " + string(e) }

func (s *Session) decodeCmd(data []byte) (instruction, error) {
	const clcNop = 1
	const clcDisconnect = 2
	const clcMove = 3
	const clcStringCmd = 4
	switch data[0] {
	case clcNop:
	case clcDisconnect:
		return func() error {
			s.Remove()
			return nil
		}, nil
	case clcMove:
		return func() error {
			return s.Move(data)
		}, nil
	case clcStringCmd:
		return func() error {
			return s.StringCmd(data)
		}, nil
	}
	return nil, errInvalidInstruction("none")

}

func (s *Session) Move(data []byte) error {
	var datum struct {
		Ping    float32
		Angle   [3]int8
		Move    [3]int16 // forward, side, up
		Button  int8
		Impulse int8
	}
	return binary.Read(bytes.NewReader(data), binary.LittleEndian, &datum)
}

func (s *Session) StringCmd(data []byte) error {
	i := 0
	for data[i] != 0 && data[i] != 255 {
		i++
	}
	log.Println("command:", string(data[:i]))
	return nil
}

const clientDisconnect = 2

var errStaleDatagram = errors.New("stale datagram")

func (s *Session) handleUnreliable(pb *datagram) error {
	if pb.Before(s.Seq) {
		return errStaleDatagram
	}
	if pb.After(s.Seq) {
		// XXX: Handle missed datagrams
	}
	s.Seq += pb.Seq() + 1
	inst, err := s.decodeCmd(pb.Data())
	if err != nil {
		return err
	}
	s.Buf.Add(inst)
	return nil
}

func (s *Session) loopNet(ctx context.Context) error {
	defer close(s.disconnectSig)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// XXX: DATA RACE
			dl := time.Duration(cvNetMessageTimeout.Get()) *
				time.Second
			ctx, _ := context.WithTimeout(ctx, dl)
			if err := s.loopNetCycle(ctx); err != nil {
				return err
			}
		}
	}
	return nil
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	terr, ok := err.(net.Error)
	if !ok {
		return false
	}
	return terr.Timeout()
}

func (s *Session) loopNetCycle(ctx context.Context) error {
	if dl, ok := ctx.Deadline(); ok {
		if err := s.Conn.SetReadDeadline(dl); err != nil {
			return err
		}
	}
	var data [maxDatagram]byte
	read, err := readDatagram(s.Conn, data[0:0])
	if err != nil {
		if isTimeout(err) {
			s.Remove()
			return fmt.Errorf("time out: %s", s)
		}
		return err
	}
	pbuf, err := decodePacketBuf(read)
	fmt.Println(pbuf, err)
	if err != nil {
		return err
	}
	if pbuf.IsNetCtrl() {
		return nil
	}
	if pbuf.IsUnreliable() {
		switch err := s.handleUnreliable(pbuf); err {
		case errStaleDatagram:
			return nil
		case nil:
		default:
			return err
		}
	}
	return nil
}
