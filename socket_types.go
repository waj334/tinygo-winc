package winc

import "C"
import (
	"bytes"
	"encoding/binary"
)

type SocketAddress struct {
	Family    uint16
	Port      uint16
	IPAddress uint32
	// 8 bytes
}

func (s *SocketAddress) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 8))

	binary.Write(buf, binary.LittleEndian, s.Family)
	binary.Write(buf, binary.LittleEndian, s.Port)
	binary.Write(buf, binary.LittleEndian, s.IPAddress)

	return buf.Bytes()
}

func (s *SocketAddress) read(buf []byte) {
	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.LittleEndian, &s.Family)
	binary.Read(reader, binary.LittleEndian, &s.Port)
	binary.Read(reader, binary.LittleEndian, &s.IPAddress)
}

type AcceptReply struct {
	Address       SocketAddress
	ListenSock    int8
	ConnectedSock int8
	AppDataOffset uint16

	// 12 bytes
}

func (a *AcceptReply) read(buf []byte) {
	a.Address.read(buf[:8])
	reader := bytes.NewReader(buf[8:])
	binary.Read(reader, binary.LittleEndian, &a.ListenSock)
	binary.Read(reader, binary.LittleEndian, &a.ConnectedSock)
	binary.Read(reader, binary.LittleEndian, &a.AppDataOffset)
}

type bindCmd struct {
	address SocketAddress
	socket  int8
	/* padding byte */
	sessionID uint16

	// 12 bytes
}

func (c *bindCmd) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 12))

	buf.Write(c.address.bytes())
	buf.WriteByte(byte(c.socket))
	buf.WriteByte(0)
	binary.Write(buf, binary.LittleEndian, c.sessionID)

	return buf.Bytes()
}

type BindReply struct {
	Socket    int8
	Status    int8
	SessionID uint16

	// 4 bytes
}

func (b *BindReply) read(buf []byte) {
	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.LittleEndian, &b.Socket)
	binary.Read(reader, binary.LittleEndian, &b.Status)
	binary.Read(reader, binary.LittleEndian, &b.SessionID)
}

type CloseCmd struct {
	socket int8
	/* padding byte */
	sessionID uint16

	// 4 bytes
}

func (c *CloseCmd) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 4))

	binary.Write(buf, binary.LittleEndian, c.socket)
	buf.WriteByte(0)
	binary.Write(buf, binary.LittleEndian, c.sessionID)

	return buf.Bytes()
}

type connectCmd struct {
	address   SocketAddress
	socket    int8
	sslFlags  byte
	sessionID uint16

	// 12 bytes
}

func (c *connectCmd) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 12))

	buf.Write(c.address.bytes())
	buf.WriteByte(byte(c.socket))
	buf.WriteByte(c.sslFlags)
	binary.Write(buf, binary.LittleEndian, c.sessionID)

	return buf.Bytes()
}

type ConnectReply struct {
	Socket    int8
	Error     int8
	ExtraData uint16

	// 4 bytes
}

func (c *ConnectReply) read(buf []byte) {
	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.LittleEndian, &c.Socket)
	binary.Read(reader, binary.LittleEndian, &c.Error)
	binary.Read(reader, binary.LittleEndian, &c.ExtraData)
}

type dnsReply struct {
	hostName [64]byte
	hostIP   uint32

	// 68 bytes
}

func (d *dnsReply) read(buf []byte) {
	reader := bytes.NewReader(buf)
	reader.Read(d.hostName[:])
	binary.Read(reader, binary.LittleEndian, &d.hostIP)
}

type listenCmd struct {
	socket    int8
	backlog   byte
	sessionID uint16

	// 4 bytes
}

func (c *listenCmd) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 4))

	buf.WriteByte(byte(c.socket))
	buf.WriteByte(c.backlog)
	binary.Write(buf, binary.LittleEndian, c.sessionID)

	return buf.Bytes()
}

type ListenReply struct {
	Socket    int8
	Status    int8
	SessionID uint16

	// 4 bytes
}

func (l *ListenReply) read(buf []byte) {
	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.LittleEndian, &l.Socket)
	binary.Read(reader, binary.LittleEndian, &l.Status)
	binary.Read(reader, binary.LittleEndian, &l.SessionID)
}

type recvCmd struct {
	timeout uint32
	socket  int8
	/* padding byte */
	sessionID uint16
	bufLen    uint16

	// 10 bytes
}

func (r *recvCmd) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 10))

	binary.Write(buf, binary.LittleEndian, r.timeout)
	binary.Write(buf, binary.LittleEndian, r.socket)
	buf.WriteByte(0)
	binary.Write(buf, binary.LittleEndian, r.sessionID)
	binary.Write(buf, binary.LittleEndian, r.bufLen)

	return buf.Bytes()
}

type RecvReply struct {
	RemoteAddress SocketAddress
	RecvStatus    int16
	DataOffset    uint16
	Socket        int8
	/* padding byte */
	SessionID uint16

	// 16 bytes
}

func (r *RecvReply) read(buf []byte) {
	r.RemoteAddress.read(buf[:8])

	reader := bytes.NewReader(buf[8:])
	binary.Read(reader, binary.LittleEndian, &r.RecvStatus)
	binary.Read(reader, binary.LittleEndian, &r.DataOffset)
	binary.Read(reader, binary.LittleEndian, &r.Socket)
	reader.ReadByte()
	binary.Read(reader, binary.LittleEndian, &r.SessionID)
}

type sendCmd struct {
	socket int8
	/* padding byte */
	dataSize  uint16
	address   SocketAddress
	sessionID uint16
	/* padding [2]byte */

	// 16 bytes
}

func (s *sendCmd) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 16))

	binary.Write(buf, binary.LittleEndian, s.socket)
	buf.WriteByte(0)
	binary.Write(buf, binary.LittleEndian, s.dataSize)
	buf.Write(s.address.bytes())
	binary.Write(buf, binary.LittleEndian, s.sessionID)
	buf.WriteByte(0)
	buf.WriteByte(0)

	return buf.Bytes()
}

type SendReply struct {
	Socket int8
	/* padding byte */
	SentBytes int16
	SessionID uint16
	/* padding [2]byte */

	// 8 bytes
}

func (s *SendReply) read(buf []byte) {
	reader := bytes.NewReader(buf)
	binary.Read(reader, binary.LittleEndian, &s.Socket)
	reader.ReadByte()
	binary.Read(reader, binary.LittleEndian, &s.SentBytes)
	binary.Read(reader, binary.LittleEndian, &s.SessionID)
}

type setSocketOptCmd struct {
	optionValue uint32
	socket      int8
	option      byte
	sessionID   uint16

	// 8 bytes
}

func (s *setSocketOptCmd) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 8))

	binary.Write(buf, binary.LittleEndian, s.optionValue)
	binary.Write(buf, binary.LittleEndian, s.socket)
	buf.WriteByte(s.option)
	binary.Write(buf, binary.LittleEndian, s.sessionID)

	return buf.Bytes()
}

type sslSetSockOptCmd struct {
	socket    int8
	option    byte
	sessionID uint16
	optLen    uint32
	optVal    [64]byte

	// 72 bytes
}

func (s *sslSetSockOptCmd) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 72))

	binary.Write(buf, binary.LittleEndian, s.socket)
	buf.WriteByte(s.option)
	binary.Write(buf, binary.LittleEndian, s.sessionID)
	binary.Write(buf, binary.LittleEndian, s.optLen)
	buf.Write(s.optVal[:])

	return buf.Bytes()
}

type sslSocketCreateCmd struct {
	socket int8
	/* padding [3]byte */

	// 4 bytes
}

func (s *sslSocketCreateCmd) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 4))

	binary.Write(buf, binary.LittleEndian, s.socket)
	buf.WriteByte(0)
	buf.WriteByte(0)
	buf.WriteByte(0)

	return buf.Bytes()
}
