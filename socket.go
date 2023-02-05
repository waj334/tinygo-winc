/*
MIT License

Copyright (c) 2022 waj334

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package winc

import (
	"net"
	"sync"
	"time"

	"github.com/waj334/tinygo-winc/protocol"
	"github.com/waj334/tinygo-winc/protocol/types"
)

type (
	SocketType   uint8
	SocketConfig uint8

	Socket struct {
		sockfd          int8
		connectedSockfd int8

		offset    uint16
		sessionId uint16
		sslFlags  uint8

		driver *WINC
		mutex  sync.Mutex

		bufferAddr uint32
		bufferLen  int

		// Socket reply channels
		acceptChan  chan int8
		bindChan    chan *types.BindReply
		connectChan chan *types.ConnectReply
		listenChan  chan *types.ListenReply
		recvChan    chan *types.RecvReply
		sendChan    chan *types.SendReply

		// Misc channels
		dnsChan chan *types.DnsReply

		recvDeadline time.Time
		sendDeadline time.Time

		addr net.Addr
	}

	Sockaddr struct {
		Port    uint16
		Address uint32
	}
)

const (
	OpcodeSocketInvalid         protocol.OpcodeId = 0x00
	OpcodeSocketBind            protocol.OpcodeId = 0x41
	OpcodeSocketListen          protocol.OpcodeId = 0x42
	OpcodeSocketAccept          protocol.OpcodeId = 0x43
	OpcodeSocketConnect         protocol.OpcodeId = 0x44
	OpcodeSocketSend            protocol.OpcodeId = 0x45
	OpcodeSocketRecv            protocol.OpcodeId = 0x46
	OpcodeSocketSendTo          protocol.OpcodeId = 0x47
	OpcodeSocketRecvFrom        protocol.OpcodeId = 0x48
	OpcodeSocketClose           protocol.OpcodeId = 0x49
	OpcodeSocketDnsResolve      protocol.OpcodeId = 0x4A
	OpcodeSocketSslConnect      protocol.OpcodeId = 0x4B
	OpcodeSocketSslSend         protocol.OpcodeId = 0x4C
	OpcodeSocketSslRecv         protocol.OpcodeId = 0x4D
	OpcodeSocketSslClose        protocol.OpcodeId = 0x4E
	OpcodeSocketSetSocketOption protocol.OpcodeId = 0x4F
	OpcodeSocketSslCreate       protocol.OpcodeId = 0x50
	OpcodeSocketSslSetSockOpt   protocol.OpcodeId = 0x51
	OpcodeSocketPing            protocol.OpcodeId = 0x52
	OpcodeSocketSslSetCsList    protocol.OpcodeId = 0x53
	OpcodeSocketSslBind         protocol.OpcodeId = 0x54
	OpcodeSocketSslExpCheck     protocol.OpcodeId = 0x55
	OpcodeSocketSecure          protocol.OpcodeId = 0x56
	OpcodeSocketSslConnectAlpn  protocol.OpcodeId = 0x57

	_NBIT31 = 0x80000000
	_NBIT30 = 0x40000000
	_NBIT29 = 0x20000000
	_NBIT28 = 0x10000000
	_NBIT27 = 0x08000000
	_NBIT26 = 0x04000000
	_NBIT25 = 0x02000000
	_NBIT24 = 0x01000000
	_NBIT23 = 0x00800000
	_NBIT22 = 0x00400000
	_NBIT21 = 0x00200000
	_NBIT20 = 0x00100000
	_NBIT19 = 0x00080000
	_NBIT18 = 0x00040000
	_NBIT17 = 0x00020000
	_NBIT16 = 0x00010000
	_NBIT15 = 0x00008000
	_NBIT14 = 0x00004000
	_NBIT13 = 0x00002000
	_NBIT12 = 0x00001000
	_NBIT11 = 0x00000800
	_NBIT10 = 0x00000400
	_NBIT9  = 0x00000200
	_NBIT8  = 0x00000100
	_NBIT7  = 0x00000080
	_NBIT6  = 0x00000040
	_NBIT5  = 0x00000020
	_NBIT4  = 0x00000010
	_NBIT3  = 0x00000008
	_NBIT2  = 0x00000004
	_NBIT1  = 0x00000002
	_NBIT0  = 0x00000001

	tlsRecordHeaderLength uint16 = 5
	ethernetHeaderOffset  uint16 = 34
	ethernetHeaderLength  uint16 = 14
	tcpIpHeaderLength     uint16 = 40
	udpIpHeaderLength     uint16 = 28

	ipPacketOffset = ethernetHeaderLength + ethernetHeaderOffset - 8

	tcpTxPacketOffset = ipPacketOffset + tcpIpHeaderLength
	udpTxPacketOffset = ipPacketOffset + udpIpHeaderLength
	sslTxPacketOffset = tcpTxPacketOffset + tlsRecordHeaderLength

	sslFlagsActive       = _NBIT0
	sslFlagsBypassX509   = _NBIT1
	sslFlags2Reserved    = _NBIT2
	sslFlags3Reserved    = _NBIT3
	sslFlagsCacheSession = _NBIT4
	sslFlagsNoTxCopy     = _NBIT5
	sslFlagsCheckSni     = _NBIT6
	sslFlagsDelay        = _NBIT7

	hostnameMaxLength = 64
	maxTcpSocket      = 7
	maxUdpSocket      = 4
	maxSocket         = maxTcpSocket + maxUdpSocket
	mtu               = 256

	SocketTypeStream   SocketType = 1
	SocketTypeDatagram SocketType = 2

	SocketConfigSslOff   SocketConfig = 0
	SocketConfigSslOn    SocketConfig = 1
	SocketConfigSslDelay SocketConfig = 2

	SolSocket    = 1
	SolSslSocket = 2

	afInet uint16 = 2
)

func (s *Socket) Listen(backlog int) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sockfd < 0 {
		return ErrSocketInvalid
	}

	// Get unique session id

	strListen := types.ListenCmd{
		Sock:         s.sockfd,
		U8BackLog:    uint8(backlog),
		U16SessionID: s.sessionId,
	}

	// Send the request to the device
	if err = s.driver.hif.Send(GroupIP, OpcodeSocketListen, strListen.Bytes(), nil, 0); err != nil {
		return
	}

	// Wait for WINC to accept the incoming connection
	var strListenReply *types.ListenReply
	select {
	case strListenReply = <-s.listenChan:
	}

	if strListenReply.S8Status < 0 {
		// Return the error
		return SocketError(strListenReply.S8Status)
	}

	return
}

func (s *Socket) Bind(addr net.Addr) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sockfd < 0 {
		return ErrSocketInvalid
	}

	strBind := types.BindCmd{
		StrAddr: types.SockAddr{
			U16Family: afInet,
		},
		Sock:         s.sockfd,
		U16SessionID: s.sessionId,
	}

	// addr can be TCPAddr or UDPAddr
	switch actualAddr := addr.(type) {
	case *TCPAddr:
		strBind.StrAddr.U32IPAddr = actualAddr.U32IPAddr
		strBind.StrAddr.U16Port = Htons(actualAddr.U16Port)
	case *UDPAddr:
		strBind.StrAddr.U32IPAddr = actualAddr.U32IPAddr
		strBind.StrAddr.U16Port = Htons(actualAddr.U16Port)
	default:
		return ErrInvalidParameter
	}

	cmd := OpcodeSocketBind
	if s.sslFlags&sslFlagsActive != 0 {
		cmd = OpcodeSocketSslBind
	}

	if err = s.driver.hif.Send(GroupIP, cmd, strBind.Bytes(), nil, 0); err != nil {
		return
	}

	var strBindReply *types.BindReply
	select {
	case strBindReply = <-s.bindChan:
	}

	if strBindReply.S8Status < 0 {
		// Return the error
		return SocketError(strBindReply.S8Status)
	}

	s.addr = addr

	return
}

func (s *Socket) Connect(addr net.Addr) (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sockfd < 0 {
		return ErrSocketInvalid
	}

	strConnect := types.ConnectCmd{
		StrAddr: types.SockAddr{
			U16Family: afInet,
		},
		Sock: s.sockfd,

		// Using the sessionId of the socket since the connect reply does not return this session id
		U16SessionID: s.sessionId,
	}

	// addr can be TCPAddr or UDPAddr
	switch actualAddr := addr.(type) {
	case *TCPAddr:
		strConnect.StrAddr.U32IPAddr = actualAddr.U32IPAddr
		strConnect.StrAddr.U16Port = Htons(actualAddr.U16Port)
	case *UDPAddr:
		strConnect.StrAddr.U32IPAddr = actualAddr.U32IPAddr
		strConnect.StrAddr.U16Port = Htons(actualAddr.U16Port)
	default:
		return ErrInvalidParameter
	}

	cmd := OpcodeSocketConnect
	if s.sslFlags&sslFlagsActive != 0 {
		cmd = OpcodeSocketSslConnect
		strConnect.U8SslFlags = s.sslFlags
	}

	if err = s.driver.hif.Send(GroupIP, cmd, strConnect.Bytes(), nil, 0); err != nil {
		return
	}

	var strConnectReply *types.ConnectReply

	// Wait for the response
	select {
	case strConnectReply = <-s.connectChan:
	}

	if strConnectReply.S8Error != 0 {
		return SocketError(strConnectReply.S8Error)
	}

	// NOTE: Extra data is the u16AppDataOffset member of the union in the original tstrConnectReply struct
	s.offset = strConnectReply.U16ExtraData - protocol.M2M_HIF_HDR_OFFSET

	// Keep the address of the remote connection
	s.addr = addr

	return
}

func (s *Socket) Shutdown() (err error) {
	if s.sockfd < 0 {
		return ErrSocketInvalid
	}

	cmd := OpcodeSocketClose
	if s.sslFlags&sslFlagsActive != 0 {
		cmd = OpcodeSocketSslClose
	}

	strClose := types.CloseCmd{
		Sock:         s.sockfd,
		U16SessionID: s.sessionId,
	}

	if err = s.driver.hif.Send(GroupIP, cmd, strClose.Bytes(), nil, 0); err != nil {
		return
	}

	// Garbage collect later
	s.driver.sockets[s.sockfd] = nil

	// Invalidate the socket so no further request can be made
	s.sockfd = -1
	return
}

func (s *Socket) Secure() (err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sockfd < 0 {
		return ErrSocketInvalid
	}

	flags := s.sslFlags
	if flags&sslFlagsDelay == 0 || flags&sslFlagsActive == 0 || s.offset == 0 {
		return ErrSocketInvalidArg
	}

	s.sslFlags &= sslFlagsDelay
	strConnect := types.ConnectCmd{
		Sock:         s.sockfd,
		U8SslFlags:   s.sslFlags,
		U16SessionID: s.sessionId,
	}

	if err = s.driver.hif.Send(GroupIP, OpcodeSocketSecure, strConnect.Bytes(), nil, 0); err != nil {
		return
	}

	// Wait for the response
	var strConnectReply *types.ConnectReply
	select {
	case strConnectReply = <-s.connectChan:
	}

	if strConnectReply.S8Error < 0 {
		return SocketError(strConnectReply.S8Error)
	}

	// NOTE: Extra data is the u16AppDataOffset member of the union in the original tstrConnectReply struct
	s.offset = strConnectReply.U16ExtraData - protocol.M2M_HIF_HDR_OFFSET

	return
}

func (s *Socket) Setsockopt(level, name int, value []byte) (err error) {
	if s.sockfd < 0 {
		return ErrSocketInvalid
	}

	cmd := OpcodeSocketSetSocketOption
	var control []byte
	if level == SolSslSocket {
		cmd = OpcodeSocketSslSetSockOpt
		var sslFlag int
		if len(value) == 4 {
			switch name {
			case types.SO_SSL_BYPASS_X509_VERIF:
				sslFlag = sslFlagsBypassX509
			case types.SO_SSL_ENABLE_SESSION_CACHING:
				sslFlag = sslFlagsCacheSession
			case types.SO_SSL_ENABLE_SNI_VALIDATION:
				sslFlag = sslFlagsCheckSni
			}
		}

		if sslFlag != 0 {
			optVal := protocol.ToUint32(value)
			if optVal != 0 {
				s.sslFlags |= uint8(sslFlag)
			} else {
				s.sslFlags &= ^uint8(sslFlag)
			}
			return
		} else if ((name == types.SO_SSL_SNI) && (len(value) < 64)) || ((name == types.SO_SSL_ALPN) && (len(value) <= 32)) {
			strSslSetSockOpt := types.SSLSetSockOptCmd{
				Sock:         s.sockfd,
				U8Option:     uint8(name),
				U16SessionID: s.sessionId,
			}

			copy(strSslSetSockOpt.Au8OptVal[:], value)
			control = strSslSetSockOpt.Bytes()
		} else {
			return ErrSocketInvalidArg
		}
	} else if level == SolSocket && len(value) == 4 {
		strSetSockOpt := types.SetSocketOptCmd{
			U32OptionValue: protocol.ToUint32(value),
			Sock:           s.sockfd,
			U8Option:       uint8(name),
			U16SessionID:   s.sessionId,
		}

		control = strSetSockOpt.Bytes()
	}

	if err = s.driver.hif.Send(GroupIP, cmd|protocol.OpcodeReqDataPkt, control, nil, 0); err != nil {
		return ErrSocketInvalid
	}

	return
}

func (s *Socket) Send(buf []byte, deadline time.Time) (sz int, err error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sockfd < 0 {
		return 0, ErrSocketInvalid
	}

	cmd := OpcodeSocketSend
	if s.sslFlags&sslFlagsActive != 0 && s.sslFlags&sslFlagsDelay == 0 {
		cmd = OpcodeSocketSslSend
	}

	strSend := types.SendCmd{
		Sock:         s.sockfd,
		U16DataSize:  uint16(len(buf)),
		U16SessionID: s.sessionId,
	}

	if err = s.driver.hif.Send(GroupIP, cmd|protocol.OpcodeReqDataPkt, strSend.Bytes(), buf, s.offset); err != nil {
		return 0, ErrSocketBufferFull
	}

	// Wait for the response
	var strSendReply *types.SendReply
	select {
	case strSendReply = <-s.sendChan:
	}

	// Check for error
	if strSendReply.S16SentBytes < 0 {
		return 0, SocketError(strSendReply.S16SentBytes)
	}

	sz = int(strSendReply.S16SentBytes)
	return
}

func (s *Socket) SendTo(buf []byte, addr net.Addr, deadline time.Time) (sz int, err error) {
	if s.sockfd < 0 {
		return 0, ErrSocketInvalid
	}

	strSend := types.SendCmd{
		Sock:        s.sockfd,
		U16DataSize: uint16(len(buf)),
		StrAddr: types.SockAddr{
			U16Family: afInet,
		},
		U16SessionID: s.sessionId,
	}

	// addr can be TCPAddr or UDPAddr
	switch actualAddr := addr.(type) {
	case *TCPAddr:
		strSend.StrAddr.U32IPAddr = actualAddr.U32IPAddr
		strSend.StrAddr.U16Port = Htons(actualAddr.U16Port)
	case *UDPAddr:
		strSend.StrAddr.U32IPAddr = actualAddr.U32IPAddr
		strSend.StrAddr.U16Port = Htons(actualAddr.U16Port)
	default:
		return 0, ErrInvalidParameter
	}

	if err = s.driver.hif.Send(GroupIP, OpcodeSocketSendTo|protocol.OpcodeReqDataPkt, strSend.Bytes(), buf, udpTxPacketOffset); err != nil {
		return 0, ErrSocketBufferFull
	}

	// Wait for the response
	var strSendReply *types.SendReply
	select {
	case strSendReply = <-s.sendChan:
	}

	// Check for error
	if strSendReply.S16SentBytes < 0 {
		return 0, SocketError(strSendReply.S16SentBytes)
	}

	sz = int(strSendReply.S16SentBytes)
	return
}

func (s *Socket) Recv(buf []byte, deadline time.Time) (sz int, err error) {
	// Block concurrent reads
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sockfd < 0 {
		return 0, ErrSocketInvalid
	}

	// Request that the firmware receives more data
	cmd := OpcodeSocketRecv
	if s.sslFlags&sslFlagsActive != 0 && s.sslFlags&sslFlagsDelay == 0 {
		cmd = OpcodeSocketSslRecv
	}

	var timeout int64
	if !deadline.IsZero() {
		timeout = deadline.Sub(time.Now()).Milliseconds()
		if timeout <= 0 {
			return sz, ErrSocketTimeout
		}
	} else {
		timeout = 0xFFFFFFFF
	}

	strRecv := types.RecvCmd{
		U32Timeoutmsec: uint32(timeout),
		Sock:           s.sockfd,
		U16BufLen:      uint16(len(buf)),
		U16SessionID:   s.driver.getSessionId(),
	}

	if err = s.driver.hif.Send(GroupIP, cmd, strRecv.Bytes(), nil, 0); err != nil {
		return 0, ErrSocketBufferFull
	}

	// Wait for the reply
	var strRecvReply *types.RecvReply
	select {
	case strRecvReply = <-s.recvChan:
		sz = int(strRecvReply.S16RecvStatus)
	}

	if strRecvReply.S16RecvStatus < 0 {
		// Return the amount of data actually read and the error
		return sz, SocketError(strRecvReply.S16RecvStatus)
	}

	// Receive the payload
	if err = s.driver.hif.Receive(s.bufferAddr, buf[:sz], true); err != nil {
		return -14, err
	}

	return
}

func (s *Socket) RecvFrom(buf []byte, addr net.Addr, deadline time.Time) (sz int, err error) {
	return s.Recv(buf, deadline)
}

func (w *WINC) Socket(sockType SocketType, config SocketConfig) (socket *Socket, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	sockfd := -1

	if sockType == SocketTypeStream {
		// Find available TCP socket
		for i := 0; i < maxTcpSocket; i++ {
			if w.sockets[i] == nil {
				sockfd = i
				break
			}
		}
	} else if sockType == SocketTypeDatagram {
		// Find available UDP socket
		for i := maxTcpSocket; i < maxSocket; i++ {
			if w.sockets[i] == nil {
				sockfd = i
				break
			}
		}
	}

	if sockfd >= 0 {
		socket = &Socket{
			sockfd: int8(sockfd),
			driver: w,

			acceptChan:  make(chan int8, 1),
			bindChan:    make(chan *types.BindReply, 1),
			connectChan: make(chan *types.ConnectReply, 1),
			listenChan:  make(chan *types.ListenReply, 1),
			recvChan:    make(chan *types.RecvReply, 1),
			sendChan:    make(chan *types.SendReply, 1),
		}

		if sockType == SocketTypeStream && config != SocketConfigSslOff {
			// Create TLS enabled socket
			strSSLCreate := types.SSLSocketCreateCmd{
				SslSock: int8(sockfd),
			}

			if err = w.hif.Send(GroupWIFI, OpcodeSocketSslCreate, strSSLCreate.Bytes(), nil, 0); err != nil {
				return nil, err
			}

			// Set TLS flags
			socket.sslFlags = sslFlagsActive | sslFlagsNoTxCopy
			if config == SocketConfigSslDelay {
				socket.sslFlags |= sslFlagsDelay
			}
		}

		// Get unique session id
		socket.sessionId = w.getSessionId()
		w.sockets[sockfd] = socket
	} else {
		err = ErrNoAvailableSocket
	}

	return
}

// SocketByDescriptor returns the pointer to an existing socket by it file descriptor. The function is useful for when
// the driver accepts an incoming connection since it will automatically open a socket for it.
func (w *WINC) SocketByDescriptor(sockfd int8) (*Socket, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if sockfd < 0 || sockfd >= maxSocket {
		return nil, ErrSocketDoesNotExist
	}

	return w.sockets[sockfd], nil
}

func (w *WINC) GetHostByName(hostname string) (address uint32, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	buf := make([]byte, len(hostname)+1)
	if len(hostname) <= hostnameMaxLength {
		copy(buf, hostname)
		if err = w.hif.Send(GroupIP, OpcodeSocketDnsResolve, buf, nil, 0); err != nil {
			return
		}

		select {
		case reply := <-w.callbackChan:
			strDnsReply := reply.(*types.DnsReply)
			address = strDnsReply.U32HostIP
		}
	}

	return
}

// getSessionId returns a unique id number. This call is thread-safe.
func (w *WINC) getSessionId() (id uint16) {
	w.sessionCounterMutex.Lock()
	id = w.sessionCounter
	w.sessionCounter++
	w.sessionCounterMutex.Unlock()

	return
}

func (w *WINC) socketCallback(id protocol.OpcodeId, sz uint16, address uint32) (data any, err error) {
	switch id {
	case OpcodeSocketAccept:
		strAcceptReply := &types.AcceptReply{}
		if err = w.hif.Receive(address, strAcceptReply.Bytes(), false); err != nil {
			return
		}

		strAcceptReply.Deref()
		strAcceptReply.Free()

		if strAcceptReply.SConnectedSock > 0 {
			// Create a socket struct for the connected socket
			w.sockets[strAcceptReply.SConnectedSock] = &Socket{
				sockfd:    strAcceptReply.SConnectedSock,
				sslFlags:  w.sockets[strAcceptReply.SListenSock].sslFlags,
				sessionId: w.getSessionId(),
				offset:    strAcceptReply.U16AppDataOffset - protocol.M2M_HIF_HDR_OFFSET,
				driver:    w,
			}

			switch w.sockets[strAcceptReply.SListenSock].addr.(type) {
			case *TCPAddr:
				w.sockets[strAcceptReply.SConnectedSock].addr = &TCPAddr{
					U16Family: afInet,
					U16Port:   strAcceptReply.StrAddr.U16Port,
					U32IPAddr: strAcceptReply.StrAddr.U32IPAddr,
				}
			case *UDPAddr:
				w.sockets[strAcceptReply.SConnectedSock].addr = &UDPAddr{
					U16Family: afInet,
					U16Port:   strAcceptReply.StrAddr.U16Port,
					U32IPAddr: strAcceptReply.StrAddr.U32IPAddr,
				}
			}

			// Signal that a socket is ready
			w.sockets[strAcceptReply.SListenSock].acceptChan <- strAcceptReply.SConnectedSock
		} else {
			w.sockets[strAcceptReply.SListenSock].acceptChan <- -1
		}

		data = strAcceptReply
	case OpcodeSocketBind, OpcodeSocketSslBind:
		strBindReply := &types.BindReply{}
		if err = w.hif.Receive(address, strBindReply.Bytes(), false); err != nil {
			return
		}

		strBindReply.Deref()
		strBindReply.Free()

		if w.sockets[strBindReply.Sock] != nil {
			w.sockets[strBindReply.Sock].bindChan <- strBindReply
		}

		data = strBindReply
	case OpcodeSocketConnect, OpcodeSocketSslConnect:
		strConnectReply := &types.ConnectReply{}
		if err = w.hif.Receive(address, strConnectReply.Bytes(), false); err != nil {
			return
		}

		strConnectReply.Deref()
		strConnectReply.Free()

		if w.sockets[strConnectReply.Sock] != nil {
			w.sockets[strConnectReply.Sock].connectChan <- strConnectReply
		}

		data = strConnectReply
	case OpcodeSocketListen:
		strListenReply := &types.ListenReply{}
		if err = w.hif.Receive(address, strListenReply.Bytes(), false); err != nil {
			return
		}

		strListenReply.Deref()
		strListenReply.Free()

		if w.sockets[strListenReply.Sock] != nil {
			w.sockets[strListenReply.Sock].listenChan <- strListenReply
		}

		data = strListenReply
	case OpcodeSocketRecv, OpcodeSocketSslRecv, OpcodeSocketRecvFrom:
		strRecvReply := &types.RecvReply{}
		if err = w.hif.Receive(address, strRecvReply.Bytes(), false); err != nil {
			return
		}

		strRecvReply.Deref()
		strRecvReply.Free()
		if strRecvReply.Sock >= 0 && strRecvReply.Sock < maxSocket {
			if w.sockets[strRecvReply.Sock] != nil {
				if strRecvReply.S16RecvStatus > 0 && strRecvReply.S16RecvStatus < int16(sz) {
					// Cache data location for Recv can retrieve the data from the WINC firmware directly
					w.sockets[strRecvReply.Sock].bufferAddr = address + uint32(strRecvReply.U16DataOffset)
				}

				w.sockets[strRecvReply.Sock].recvChan <- strRecvReply
			}

			data = strRecvReply
		} else {
			return nil, ErrSocketDoesNotExist
		}
	case OpcodeSocketSend, OpcodeSocketSslSend, OpcodeSocketSendTo:
		strSendReply := &types.SendReply{}
		if err = w.hif.Receive(address, strSendReply.Bytes(), false); err != nil {
			return
		}

		strSendReply.Deref()
		strSendReply.Free()

		if w.sockets[strSendReply.Sock] != nil {
			w.sockets[strSendReply.Sock].sendChan <- strSendReply
		}

		data = strSendReply
	case OpcodeSocketDnsResolve:
		strDnsReply := &types.DnsReply{}
		if err = w.hif.Receive(address, strDnsReply.Bytes(), false); err != nil {
			return
		}

		strDnsReply.Deref()
		strDnsReply.Free()

		w.callbackChan <- strDnsReply
		data = strDnsReply
	}

	return
}

//go:inline
func Htons(port uint16) uint16 {
	return (port << 8) | (port >> 8)
}

//go:inline
func Htonl(m uint32) uint32 {
	return (m << 24) | ((m & 0x0000FF00) << 8) | ((m & 0x00FF0000) >> 8) | (m >> 24)
}
