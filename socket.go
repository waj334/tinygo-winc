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

		acceptChan chan int8

		bufferAddr uint32
		bufferLen  int

		callbackChan SyncMap[uint16, chan<- any]
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

	_NBIT31 uint32 = 0x80000000
	_NBIT30 uint32 = 0x40000000
	_NBIT29 uint32 = 0x20000000
	_NBIT28 uint32 = 0x10000000
	_NBIT27 uint32 = 0x08000000
	_NBIT26 uint32 = 0x04000000
	_NBIT25 uint32 = 0x02000000
	_NBIT24 uint32 = 0x01000000
	_NBIT23 uint32 = 0x00800000
	_NBIT22 uint32 = 0x00400000
	_NBIT21 uint32 = 0x00200000
	_NBIT20 uint32 = 0x00100000
	_NBIT19 uint32 = 0x00080000
	_NBIT18 uint32 = 0x00040000
	_NBIT17 uint32 = 0x00020000
	_NBIT16 uint32 = 0x00010000
	_NBIT15 uint32 = 0x00008000
	_NBIT14 uint32 = 0x00004000
	_NBIT13 uint32 = 0x00002000
	_NBIT12 uint32 = 0x00001000
	_NBIT11 uint32 = 0x00000800
	_NBIT10 uint32 = 0x00000400
	_NBIT9  uint32 = 0x00000200
	_NBIT8  uint32 = 0x00000100
	_NBIT7  uint32 = 0x00000080
	_NBIT6  uint32 = 0x00000040
	_NBIT5  uint32 = 0x00000020
	_NBIT4  uint32 = 0x00000010
	_NBIT3  uint32 = 0x00000008
	_NBIT2  uint32 = 0x00000004
	_NBIT1  uint32 = 0x00000002
	_NBIT0  uint32 = 0x00000001

	tlsRecordHeaderLength uint16 = 5
	ethernetHeaderOffset  uint16 = 34
	ethernetHeaderLength  uint16 = 14
	tcpIpHeaderLength     uint16 = 40
	udpIpHeaderLength     uint16 = 28

	ipPacketOffset = ethernetHeaderLength + ethernetHeaderOffset - 8

	tcpTxPacketOffset = ipPacketOffset + tcpIpHeaderLength
	udpTxPacketOffset = ipPacketOffset + udpIpHeaderLength
	sslTxPacketOffset = tcpTxPacketOffset + tlsRecordHeaderLength

	sslFlagsActive       = uint8(_NBIT0)
	sslFlagsBypassX509   = uint8(_NBIT1)
	sslFlags2Reserved    = uint8(_NBIT2)
	sslFlags3Reserved    = uint8(_NBIT3)
	sslFlagsCacheSession = uint8(_NBIT4)
	sslFlagsNoTxCopy     = uint8(_NBIT5)
	sslFlagsCheckSni     = uint8(_NBIT6)
	sslFlagsDelay        = uint8(_NBIT7)

	hostnameMaxLength = 64
	maxTcpSocket      = 7
	maxUdpSocket      = 4
	maxSocket         = maxTcpSocket + maxUdpSocket
	mtu               = 256

	SocketTypeStream   SocketType = 1
	SocketTypeDatagram SocketType = 2

	SocketConfigSslOff   = 0
	SocketConfigSslOn    = 1
	SocketConfigSslDelay = 2

	SolSocket    = 1
	SolSslSocket = 2

	afInet uint16 = 2
)

func (s *Socket) Listen(backlog int) (err error) {
	if s.sockfd < 0 {
		return ErrSocketInvalid
	}

	// Get unique session id

	strListen := types.ListenCmd{
		Sock:         s.sockfd,
		U8BackLog:    uint8(backlog),
		U16SessionID: s.driver.getSessionId(),
	}

	// Create a channel to receive the reply on
	replyChan := make(chan any, 1)
	s.callbackChan.Store(strListen.U16SessionID, replyChan)
	defer close(replyChan)

	// Send the request to the device
	if err = s.driver.hif.Send(GroupIP, OpcodeSocketListen, strListen.Bytes(), nil, 0); err != nil {
		return
	}

	s.mutex.Lock()

	// Wait for WINC to accept the incoming connection
	var strListenReply types.ListenReply
	select {
	case reply := <-replyChan:
		strListenReply = reply.(types.ListenReply)
	}

	if strListenReply.S8Status < 0 {
		// Return the error
		return SocketError(strListenReply.S8Status)
	}

	return
}

func (s *Socket) Bind(addr net.Addr) (err error) {
	if s.sockfd < 0 {
		return ErrSocketInvalid
	}

	strBind := types.BindCmd{
		StrAddr: types.SockAddr{
			U16Family: afInet,
		},
		Sock:         s.sockfd,
		U16SessionID: s.driver.getSessionId(),
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

	// Create a channel to receive the reply on
	replyChan := make(chan any, 1)
	s.callbackChan.Store(strBind.U16SessionID, replyChan)
	defer close(replyChan)

	if err = s.driver.hif.Send(GroupIP, cmd, strBind.Bytes(), nil, 0); err != nil {
		return
	}

	var strBindReply types.BindReply
	select {
	case reply := <-replyChan:
		strBindReply = reply.(types.BindReply)
	}

	if strBindReply.S8Status < 0 {
		// Return the error
		return SocketError(strBindReply.S8Status)
	}

	s.addr = addr

	return
}

func (s *Socket) Connect(addr net.Addr) (err error) {
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

	// Create a channel to receive the reply on
	replyChan := make(chan any, 1)
	s.callbackChan.Store(strConnect.U16SessionID, replyChan)
	defer close(replyChan)

	if err = s.driver.hif.Send(GroupIP, cmd, strConnect.Bytes(), nil, 0); err != nil {
		return
	}

	var strConnectReply types.ConnectReply

	// Wait for the response
	select {
	case reply := <-replyChan:
		strConnectReply = reply.(types.ConnectReply)
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
		U16SessionID: s.driver.getSessionId(),
	}

	// Create a channel to receive the reply on
	replyChan := make(chan any, 1)
	s.callbackChan.Store(strConnect.U16SessionID, replyChan)
	defer close(replyChan)

	if err = s.driver.hif.Send(GroupIP, OpcodeSocketSecure, strConnect.Bytes(), nil, 0); err != nil {
		return
	}

	// Wait for the response
	var strConnectReply types.ConnectReply
	select {
	case reply := <-replyChan:
		strConnectReply = reply.(types.ConnectReply)
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
				sslFlag = types.SSL_FLAGS_BYPASS_X509
			case types.SO_SSL_ENABLE_SESSION_CACHING:
				sslFlag = types.SSL_FLAGS_CACHE_SESSION
			case types.SO_SSL_ENABLE_SNI_VALIDATION:
				sslFlag = types.SSL_FLAGS_CHECK_SNI
			}
		}

		if sslFlag != 0 {
			optVal := protocol.ToUint32(value)
			if optVal != 0 {
				s.sslFlags |= uint8(optVal)
			} else {
				s.sslFlags &= ^uint8(optVal)
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
	if s.sockfd < 0 {
		return 0, ErrSocketInvalid
	}

	cmd := OpcodeSocketSend
	offset := tcpTxPacketOffset

	if s.sockfd >= maxTcpSocket {
		offset = udpTxPacketOffset
	}

	if s.sslFlags&sslFlagsActive != 0 && s.sslFlags&sslFlagsDelay == 0 {
		cmd = OpcodeSocketSslSend
		offset = s.offset
	}

	strSend := types.SendCmd{
		Sock:         s.sockfd,
		U16DataSize:  uint16(len(buf)),
		U16SessionID: s.driver.getSessionId(),
	}

	// Create a channel to receive the reply on
	replyChan := make(chan any, 1)
	s.callbackChan.Store(strSend.U16SessionID, replyChan)
	defer close(replyChan)

	if err = s.driver.hif.Send(GroupIP, cmd|protocol.OpcodeReqDataPkt, strSend.Bytes(), buf, offset); err != nil {
		s.callbackChan.Delete(strSend.U16SessionID)
		return 0, ErrSocketBufferFull
	}

	var strSendReply types.SendReply

	// Wait for the response
	if s.sendDeadline.IsZero() {
		select {
		case reply := <-replyChan:
			strSendReply = reply.(types.SendReply)
		}
	} else {
		select {
		case reply := <-replyChan:
			strSendReply = reply.(types.SendReply)
		case <-time.After(s.sendDeadline.Sub(time.Now())):
			s.callbackChan.Delete(strSend.U16SessionID)
			return 0, ErrSocketTimeout
		}
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
		U16SessionID: s.driver.getSessionId(),
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

	// Create a channel to receive the reply on
	replyChan := make(chan any, 1)
	s.callbackChan.Store(strSend.U16SessionID, replyChan)
	defer close(replyChan)

	if err = s.driver.hif.Send(GroupIP, OpcodeSocketSendTo|protocol.OpcodeReqDataPkt, strSend.Bytes(), buf, udpTxPacketOffset); err != nil {
		s.callbackChan.Delete(strSend.U16SessionID)
		return 0, ErrSocketBufferFull
	}

	var strSendReply types.SendReply

	// Wait for the response
	if s.sendDeadline.IsZero() {
		select {
		case reply := <-replyChan:
			strSendReply = reply.(types.SendReply)
		}
	} else {
		select {
		case reply := <-replyChan:
			strSendReply = reply.(types.SendReply)
		case <-time.After(s.sendDeadline.Sub(time.Now())):
			s.callbackChan.Delete(strSend.U16SessionID)
			return 0, ErrSocketTimeout
		}
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

	var timeout uint32
	if !deadline.IsZero() {
		timeout = uint32(deadline.Sub(time.Now()).Milliseconds())
		if timeout <= 0 {
			timeout = 0
		}
	} else {
		timeout = 0xFFFFFFFF
	}

	// Keep track of how much data is left to be received from the WINC firmware
	remaining := len(buf)

	// Attempt to receive data repeatedly until the input buffer is entirely filled
	for remaining != 0 {
		// Send Recv request if there are no bytes available to read
		if s.bufferLen == 0 {
			// Request that the firmware receives more data
			cmd := OpcodeSocketRecv
			if s.sslFlags&sslFlagsActive != 0 && s.sslFlags&sslFlagsDelay == 0 {
				cmd = OpcodeSocketSslRecv
			}

			strRecv := types.RecvCmd{
				U32Timeoutmsec: timeout,
				Sock:           s.sockfd,
				U16SessionID:   s.driver.getSessionId(),
				U16BufLen:      uint16(len(buf)),
			}

			// Create a channel to receive the reply on
			replyChan := make(chan any, 1)
			s.callbackChan.Store(strRecv.U16SessionID, replyChan)

			if err = s.driver.hif.Send(GroupIP, cmd, strRecv.Bytes(), nil, 0); err != nil {
				close(replyChan)
				return 0, ErrSocketBufferFull
			}

			// Wait for the reply
			var strRecvReply types.RecvReply
			select {
			case reply := <-replyChan:
				strRecvReply = reply.(types.RecvReply)
				close(replyChan)
			}

			if strRecvReply.S16RecvStatus < 0 {
				// Return the amount of data actually read and the error
				return sz, SocketError(strRecvReply.S16RecvStatus)
			}
		}

		// Calculate how much data can be received
		receiveLen := remaining
		if receiveLen > s.bufferLen {
			// Cap at the number of remaining bytes available to be read. A Recv request will be made after receiving
			// the bytes from the WINC firmware.
			receiveLen = s.bufferLen
		}

		// TODO: hif.Receive should probably report the number of bytes that were actually returned from the firmware.
		err = s.driver.hif.Receive(s.bufferAddr, buf[sz:sz+receiveLen], true)
		s.bufferLen -= receiveLen
		s.bufferAddr += uint32(receiveLen)

		sz += receiveLen
		remaining -= receiveLen
	}

	return
}

func (s *Socket) RecvFrom(buf []byte, addr net.Addr, deadline time.Time) (sz int, err error) {
	// Block concurrent reads
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.sockfd < 0 {
		return 0, ErrSocketInvalid
	}

	var timeout uint32
	if !deadline.IsZero() {
		timeout = uint32(deadline.Sub(time.Now()).Milliseconds())
		if timeout <= 0 {
			timeout = 0
		}
	} else {
		timeout = 0xFFFFFFFF
	}

	// Keep track of how much data is left to be received from the WINC firmware
	remaining := len(buf)

	// Attempt to receive data repeatedly until the input buffer is entirely filled
	for remaining != 0 {
		// Send Recv request if there are no bytes available to read
		if s.bufferLen == 0 {
			strRecv := types.RecvCmd{
				U32Timeoutmsec: timeout,
				Sock:           s.sockfd,
				U16SessionID:   s.driver.getSessionId(),
				U16BufLen:      uint16(len(buf)),
			}

			// addr can be TCPAddr or UDPAddr
			/*
				switch actualAddr := addr.(type) {
				case *net.TCPAddr:
					strConnect.StrAddr.U32IPAddr = binary.BigEndian.Uint32(actualAddr.IP)
					strConnect.StrAddr.U16Port = uint16(actualAddr.Port)
				case *net.UDPAddr:
					strConnect.StrAddr.U32IPAddr = binary.BigEndian.Uint32(actualAddr.IP)
					strConnect.StrAddr.U16Port = uint16(actualAddr.Port)
				}
			*/

			// Create a channel to receive the reply on
			replyChan := make(chan any, 1)
			s.callbackChan.Store(strRecv.U16SessionID, replyChan)

			if err = s.driver.hif.Send(GroupIP, OpcodeSocketRecvFrom, strRecv.Bytes(), nil, 0); err != nil {
				close(replyChan)
				return
			}

			var strRecvReply types.RecvReply

			// Wait for reply
			select {
			case reply := <-replyChan:
				strRecvReply = reply.(types.RecvReply)
				close(replyChan)
			}

			if strRecvReply.S16RecvStatus < 0 {
				return sz, SocketError(strRecvReply.S16RecvStatus)
			}
		}

		// Calculate how much data can be received
		receiveLen := remaining
		if receiveLen > s.bufferLen {
			// Cap at the number of remaining bytes available to be read. A Recv request will be made after receiving
			// the bytes from the WINC firmware.
			receiveLen = s.bufferLen
		}

		// TODO: hif.Receive should probably report the number of bytes that were actually returned from the firmware.
		err = s.driver.hif.Receive(s.bufferAddr, buf[sz:sz+receiveLen], true)
		s.bufferLen -= receiveLen
		s.bufferAddr += uint32(receiveLen)

		sz += receiveLen
		remaining -= receiveLen
	}

	return
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
			sockfd:     int8(sockfd),
			acceptChan: make(chan int8, 1),
			driver:     w,
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
			w.sockets[sockfd].sslFlags = sslFlagsActive | sslFlagsNoTxCopy
			if config == SocketConfigSslDelay {
				w.sockets[sockfd].sslFlags |= sslFlagsDelay
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
			strDnsReply := reply.(types.DnsReply)
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
		var strAcceptReply types.AcceptReply
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
		var strBindReply types.BindReply
		if err = w.hif.Receive(address, strBindReply.Bytes(), false); err != nil {
			return
		}

		strBindReply.Deref()
		strBindReply.Free()

		if w.sockets[strBindReply.Sock] != nil {
			if replyChan, ok := w.sockets[strBindReply.Sock].callbackChan.LoadAndDelete(strBindReply.U16SessionID); ok {
				replyChan <- strBindReply
			}
		}

		data = strBindReply
	case OpcodeSocketConnect, OpcodeSocketSslConnect:
		var strConnectReply types.ConnectReply
		if err = w.hif.Receive(address, strConnectReply.Bytes(), false); err != nil {
			return
		}

		strConnectReply.Deref()
		strConnectReply.Free()

		if w.sockets[strConnectReply.Sock] != nil {
			sessionId := w.sockets[strConnectReply.Sock].sessionId
			if replyChan, ok := w.sockets[strConnectReply.Sock].callbackChan.LoadAndDelete(sessionId); ok {
				replyChan <- strConnectReply
			}
		}

		data = strConnectReply
	case OpcodeSocketListen:
		var strListenReply types.ListenReply
		if err = w.hif.Receive(address, strListenReply.Bytes(), false); err != nil {
			return
		}

		strListenReply.Deref()
		strListenReply.Free()

		if w.sockets[strListenReply.Sock] != nil {
			if replyChan, ok := w.sockets[strListenReply.Sock].callbackChan.LoadAndDelete(strListenReply.U16SessionID); ok {
				replyChan <- strListenReply
			}
		}

		data = strListenReply
	case OpcodeSocketRecv, OpcodeSocketSslRecv, OpcodeSocketRecvFrom:
		var strRecvReply types.RecvReply
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
					w.sockets[strRecvReply.Sock].bufferLen += int(strRecvReply.S16RecvStatus)
				}

				if replyChan, ok := w.sockets[strRecvReply.Sock].callbackChan.LoadAndDelete(strRecvReply.U16SessionID); ok {
					replyChan <- strRecvReply
				}
			}
			data = strRecvReply
		} else {
			return nil, ErrSocketDoesNotExist
		}
	case OpcodeSocketSend, OpcodeSocketSslSend, OpcodeSocketSendTo:
		var strSendReply types.SendReply
		if err = w.hif.Receive(address, strSendReply.Bytes(), false); err != nil {
			return
		}

		strSendReply.Deref()
		strSendReply.Free()

		if w.sockets[strSendReply.Sock] != nil {
			if replyChan, ok := w.sockets[strSendReply.Sock].callbackChan.LoadAndDelete(strSendReply.U16SessionID); ok {
				replyChan <- strSendReply
			}
		}

		data = strSendReply
	case OpcodeSocketDnsResolve:
		var strDnsReply types.DnsReply
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
