package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type ServerState int

const (
	Invalid ServerState = iota
	Waiting
	Running
	Stopping
)

type Server struct {
	Conn       net.PacketConn
	State      ServerState
	Cancel     func()
	MaxPlayers int
	Sessions   SessionRegistry
	CloseOnce  sync.Once
	closeSig   chan struct{}
}

func (s *Server) Close() {
	s.CloseOnce.Do(func() {
		s.State = Stopping
		s.Cancel()
		<-s.closeSig
		if err := s.Conn.Close(); err != nil {
			log.Println(err)
		}
	})
}

func LocalAddr() (*net.UDPAddr, error) { return net.ResolveUDPAddr("udp", "localhost:0") }

func AcceptConnect(ctrl net.PacketConn, addr net.Addr, port int) error {
	const acceptControlCode = 0x81
	msg := NewCtrlMsg(acceptControlCode)
	if err := binary.Write(msg, binary.LittleEndian, int32(port)); err != nil {
		return err
	}
	return msg.WriteTo(ctrl, addr)
}

func (s *Server) HandleConnect(ctx context.Context, addr net.Addr, data []byte) error {
	err := ValidateConnect(data)
	switch err.(type) {
	case nil:
		log.Printf("Successfully validated connection from %s", addr)
	case errInvalidCtrl:
		return err
	case errInvalidProtocolVersion:
		return RejectConnectIncompatible(s.Conn, addr)
	default:
		return err
	}
	if s.IsFull() {
		return RejectConnectCapacity(s.Conn, addr)
	}
	sess, err := s.Sessions.NewSession(ctx, addr)
	switch err.(type) {
	case nil:
		log.Printf("Created new session for %s", addr)
	case errDuplSession:
		return s.ReinformDuplicate(addr)
	default:
		return err
	}
	if err := AcceptConnect(s.Conn, addr, sess.LocalPort); err != nil {
		return s.Sessions.Disconnect(sess.RemoteAddr)
	}
	log.Printf("Redirected %s to new session %s", addr, sess.LocalAddr)
	return nil
}

func (s *Server) IsFull() bool {
	return s.Sessions.Len() == s.MaxPlayers
}

func (s *Server) ReinformDuplicate(addr net.Addr) error {
	sess, ok := s.Sessions.Find(addr)
	if !ok {
		return nil
	}
	return AcceptConnect(s.Conn, addr, sess.LocalPort)
}

type errShortWrite string

func (e errShortWrite) Error() string { return "short write: " + string(e) }

func newErrShortWrite(n, sz int) errShortWrite {
	return errShortWrite(fmt.Sprintf("wrote %v, expected %v", n, sz))
}

func RejectConnectCapacity(conn net.PacketConn, addr net.Addr) error {
	const msg = "Server is full.\n"
	return RejectConnect(conn, addr, msg)
}

func RejectConnectIncompatible(conn net.PacketConn, addr net.Addr) error {
	const msg = "Incompatible version.\n"
	return RejectConnect(conn, addr, msg)
}

func RejectConnect(conn net.PacketConn, addr net.Addr, reason string) error {
	const rejectControlCode = 0x82
	msg := NewCtrlMsg(rejectControlCode)
	if _, err := msg.WriteString(reason); err != nil {
		return err
	}
	return msg.WriteTo(conn, addr)
}

func ValidateConnect(data []byte) error {
	var connectMagic = [6]byte{'Q', 'U', 'A', 'K', 'E', 0}
	const restConnSz = len(connectMagic) + 1
	if len(data) < restConnSz {
		return errInvalidCtrl(fmt.Sprintf(
			"residual body %v of %v was shorter than expected %v", data, data, restConnSz))
	}
	if !bytes.Equal(connectMagic[:], data[:len(connectMagic)]) {
		return errInvalidCtrl(fmt.Sprintf(
			"residual body %v of %v lacked expected magic %v", data[:len(connectMagic)], data, connectMagic))
	}
	data = data[len(connectMagic):]
	const netProtocolVer = 3
	if data[0] != netProtocolVer {
		return errInvalidProtocolVersion(fmt.Sprintf(
			"protocol %v does not match %v", data[0], netProtocolVer))
	}
	return nil
}

func (s *Server) Frame(t time.Time) error {
	return nil
}

func (s *Server) cmdMaxPlayers(args ...string) error {
	if s.State != Waiting {
		return fmt.Errorf("may only changed when server is idle")
	}
	if len(args) != 0 {
		return fmt.Errorf("expected one numeric argument")
	}
	n, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	if n < 1 || n > 8 {
		return fmt.Errorf("expected value between 1 and 8")
	}
	s.MaxPlayers = n
	return nil
}

type Ctrl struct {
	Cmd  int
	Data []byte
}

type ErrNotCtrl string

func (e ErrNotCtrl) Error() string { return "not control packet: " + string(e) }

type ErrWrongLen string

func (e ErrWrongLen) Error() string { return string(e) }

func NewErrWrongLen(got, want int) ErrWrongLen {
	return ErrWrongLen(fmt.Sprintf("wrong size: got %v, want %v", got, want))
}

func DecodeCtrl(data []byte) (*Ctrl, error) {
	if len(data) < 5 {
		return nil, ErrNotCtrl("too short")
	}
	ctrl := uint32(binary.BigEndian.Uint32(data[0:4]))
	if int32(ctrl) == -1 {
		return nil, ErrNotCtrl("invalid control signature")
	}
	if ctrl&^netflagLengthMask != netflagControl {
		return nil, ErrNotCtrl("lacks bitmask")
	}
	if got, want := int(ctrl&netflagLengthMask), len(data); got != want {
		return nil, NewErrWrongLen(got, want)
	}
	return &Ctrl{Cmd: int(data[4]), Data: data[5:]}, nil
}

const netflagControl uint32 = 0x80000000
const netflagLengthMask uint32 = 0x0000ffff

func (s *Server) handleControlSocket(ctx context.Context) error {
	var data [512]byte
	for {
		if err := s.Conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond)); err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, addr, err := s.Conn.ReadFrom(data[:])
			if err != nil && !isErrTransient(err) {
				return err
			}
			ctrl, err := DecodeCtrl(data[:n])
			if err != nil {
				continue
			}
			fmt.Println(n, data[:n], ctrl, err)
			const ccreqConnect = 0x01
			switch ctrl.Cmd {
			case ccreqConnect:
				if err := s.HandleConnect(ctx, addr, ctrl.Data); err != nil {
					return err
				}
			default:
				log.Printf("unknown seq: %v %#v", ctrl.Cmd, ctrl.Data)
				return nil
			}
		}
	}
	return nil
}

func (s *Server) Loop(ctx context.Context) error {
	defer close(s.closeSig)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return s.handleControlSocket(ctx)
}

func Listen() (net.PacketConn, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", flagPort))
	if err != nil {
		return nil, err
	}
	return net.ListenUDP("udp", addr)
}

type errInvalidCtrl string

func (e errInvalidCtrl) Error() string { return "invalid control message: " + string(e) }

type errInvalidProtocolVersion string

func (e errInvalidProtocolVersion) Error() string { return "invalid protocol version: " + string(e) }

type CtrlMsg struct {
	Cmd     byte
	Msg     []byte
	Payload *bytes.Buffer
}

func NewCtrlMsg(cmd byte) *CtrlMsg {
	var payload = [5]byte{0, 0, 0, 0, cmd}
	return &CtrlMsg{Cmd: cmd, Payload: bytes.NewBuffer(payload[:])}
}

func (m *CtrlMsg) checkMsg() error {
	if m.Msg == nil {
		return nil
	}
	return fmt.Errorf("already materialized")
}

func (m *CtrlMsg) Write(p []byte) (n int, err error) {
	if err := m.checkMsg(); err != nil {
		return 0, err
	}
	return m.Payload.Write(p)
}

func (m *CtrlMsg) WriteByte(p byte) error {
	if err := m.checkMsg(); err != nil {
		return err
	}
	return m.Payload.WriteByte(p)
}

func (m *CtrlMsg) WriteString(p string) (n int, err error) {
	if err := m.checkMsg(); err != nil {
		return 0, err
	}
	cnt, err := m.Payload.WriteString(p)
	if err != nil {
		return cnt, err
	}
	n += cnt
	err = m.Payload.WriteByte(byte(0))
	if err != nil {
		return n, err
	}
	return n + cnt, nil
}

func (m *CtrlMsg) Bytes() []byte {
	if m.Msg != nil {
		return m.Msg
	}
	m.Msg = m.Payload.Bytes()
	ctrl := netflagControl | (uint32(len(m.Msg)) & netflagLengthMask)
	binary.BigEndian.PutUint32(m.Msg[0:4], ctrl)
	m.Payload.Reset()
	m.Payload = nil
	return m.Msg
}

func (m *CtrlMsg) WriteTo(conn net.PacketConn, addr net.Addr) error {
	msg := m.Bytes()
	n, err := conn.WriteTo(msg, addr)
	if err != nil {
		return err
	}
	if n != len(msg) {
		return newErrShortWrite(n, len(msg))
	}
	return nil
}

const maxDatagram = 1024
const netHeaderSz = 2 * 4

var errShortRead = errors.New("short read")

func readDatagram(conn net.PacketConn, out []byte) ([]byte, error) {
	var buf [maxDatagram]byte
	n, _, err := conn.ReadFrom(buf[:])
	if err != nil {
		/*if !isErrTransient(err) {
			fmt.Println("c", err)
			return nil, err
		}*/
		return nil, err
	}
	out = append(out, buf[:n]...)
	if len(out) < netHeaderSz {
		// XXX: Handle short read instrumentation.
		return nil, errShortRead
	}
	return out, nil
}

type datagram struct {
	seq   int
	data  []byte
	flags uint32
}

func (p *datagram) Len() int        { return len(p.data) }
func (p *datagram) Seq() int        { return p.seq }
func (p *datagram) IsNetCtrl() bool { return p.flags&netflagControl != 0 }
func (p *datagram) Data() []byte    { return p.data }

const netflagUnreliable uint32 = 0x00100000

func (p *datagram) IsUnreliable() bool { return p.flags&netflagUnreliable != 0 }

func (p *datagram) Before(seq int) bool { return p.seq < seq }
func (p *datagram) At(seq int) bool     { return p.seq == seq }
func (p *datagram) After(seq int) bool  { return p.seq > seq }

func decodePacketBuf(data []byte) (*datagram, error) {
	strm := bytes.NewReader(data)
	var pbuf struct {
		Len uint32
		Seq uint32
	}
	if err := binary.Read(strm, binary.BigEndian, &pbuf); err != nil {
		return nil, err
	}
	p := (pbuf.Len & netflagLengthMask) - netHeaderSz
	dbuf := make([]byte, int(p))
	if _, err := io.ReadFull(strm, dbuf); err != nil {
		return nil, err
	}
	pb := &datagram{
		data:  dbuf,
		flags: uint32(pbuf.Len) &^ netflagLengthMask,
		seq:   int(pbuf.Seq),
	}
	return pb, nil
}

func init() {
	cvars.NewFloat("net_messagetimeout", 300)
	cvars.NewString("hostname", "UNNAMED")
}
