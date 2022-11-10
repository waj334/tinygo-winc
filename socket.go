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

import "C"
import (
	"sync"
	"time"
	"unsafe"

	"github.com/smallnest/ringbuffer"

	"github.com/waj334/tinygo-winc/protocol"
	"github.com/waj334/tinygo-winc/protocol/types"
)

type (
	Socket       int8
	SocketType   uint8
	SocketConfig uint8
	socketStr    struct {
		buffer         *ringbuffer.RingBuffer
		offset         uint16
		sessionId      uint16
		inUse          bool
		sslFlags       uint8
		receivePending bool
		alpnStatus     uint8
		errSource      uint8
		errCode        uint8
		mutex          sync.Mutex
		callbackChan   chan any
		acceptChan     chan types.AcceptReply
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

const (
	SocketInvalid Socket = -1
)

var (
	hostnameChan = make(chan types.DnsReply, 1)

	sockets [maxSocket]socketStr

	sessionCounterMutex sync.Mutex
	sessionCounter      uint16 = 1
)

func (w *WINC) Socket(sockType SocketType, config SocketConfig) (socket Socket, err error) {
	socket = SocketInvalid

	if sockType == SocketTypeStream {
		// Find available TCP socket
		for i := 0; i < maxTcpSocket; i++ {
			if !sockets[i].inUse {
				socket = Socket(i)
				break
			}
		}
	} else if sockType == SocketTypeDatagram {
		// Find available UDP socket
		for i := maxTcpSocket; i < maxSocket; i++ {
			if !sockets[i].inUse {
				socket = Socket(i)
				break
			}
		}
	}

	if socket >= 0 {
		println("Creating socket", socket, "/(", maxTcpSocket, ",", maxSocket, ")")
		sockets[socket] = socketStr{
			buffer: ringbuffer.New(2048),
		}

		sockets[socket].mutex.Lock()
		defer sockets[socket].mutex.Unlock()

		if sockType == SocketTypeStream && config != SocketConfigSslOff {
			// Create TLS enabled socket
			strSSLCreate := types.SSLSocketCreateCmd{
				SslSock: int8(socket),
			}

			if err = hif.Send(GroupWIFI, OpcodeSocketSslCreate, strSSLCreate.Bytes(), nil, 0); err != nil {
				return SocketInvalid, err
			}

			// Set TLS flags
			sockets[socket].sslFlags = sslFlagsActive | sslFlagsNoTxCopy
			if config == SocketConfigSslDelay {
				sockets[socket].sslFlags |= sslFlagsDelay
			}
		}

		sockets[socket].inUse = true
		sockets[socket].callbackChan = make(chan interface{}, 1)

		// Get unique session id
		sessionCounterMutex.Lock()
		sockets[socket].sessionId = sessionCounter
		sessionCounter++
		sessionCounterMutex.Unlock()
	} else {
		err = ErrNoAvailableSocket
	}

	return
}

func (w *WINC) Accept(socket Socket, addr Sockaddr) (err error) {
	return ErrSocketFuncNotImplemented
}

func (w *WINC) Listen(socket Socket, backlog int) (err error) {
	sockets[socket].mutex.Lock()
	defer sockets[socket].mutex.Unlock()

	if socket < 0 || socket > maxSocket || !sockets[socket].inUse {
		return ErrSocketInvalidArg
	}

	strListen := types.ListenCmd{
		Sock:         int8(socket),
		U8BackLog:    uint8(backlog),
		U16SessionID: sockets[socket].sessionId,
	}

	if err = hif.Send(GroupIP, OpcodeSocketListen, strListen.Bytes(), nil, 0); err != nil {
		return
	}

	// Wait for WINC to accept the incoming connection
	var strListenReply types.ListenReply
	select {
	case reply := <-sockets[socket].callbackChan:
		strListenReply = reply.(types.ListenReply)
	}

	if strListenReply.S8Status < 0 {
		// Return the error
		return SocketError(strListenReply.S8Status)
	}

	return
}

func (w *WINC) Bind(socket Socket, addr Sockaddr) (err error) {
	sockets[socket].mutex.Lock()
	defer sockets[socket].mutex.Unlock()

	if socket < 0 || socket > maxSocket || !sockets[socket].inUse {
		return ErrSocketInvalidArg
	}

	strBind := types.BindCmd{
		StrAddr: types.SockAddr{
			U16Family: afInet,
			U16Port:   addr.Port,
			U32IPAddr: addr.Address,
		},
		Sock:         int8(socket),
		U16SessionID: sockets[socket].sessionId,
	}

	cmd := OpcodeSocketBind
	if sockets[socket].sslFlags&sslFlagsActive != 0 {
		cmd = OpcodeSocketSslBind
	}

	if err = hif.Send(GroupIP, cmd, strBind.Bytes(), nil, 0); err != nil {
		return
	}

	return
}

func (w *WINC) Connect(socket Socket, addr Sockaddr) (err error) {
	sockets[socket].mutex.Lock()
	defer sockets[socket].mutex.Unlock()

	if sockets[socket].inUse {
		strConnect := types.ConnectCmd{
			StrAddr: types.SockAddr{
				U16Family: afInet,
				U16Port:   addr.Port,
				U32IPAddr: addr.Address,
			},
			Sock:         int8(socket),
			U16SessionID: sockets[socket].sessionId,
		}

		cmd := OpcodeSocketConnect
		if sockets[socket].sslFlags&sslFlagsActive != 0 {
			cmd = OpcodeSocketSslConnect
			strConnect.U8SslFlags = sockets[socket].sslFlags
		}

		if err = hif.Send(GroupIP, cmd, strConnect.Bytes(), nil, 0); err != nil {
			return
		}

		var strConnectReply types.ConnectReply

		// Wait for the response
		select {
		case reply := <-sockets[socket].callbackChan:
			strConnectReply = reply.(types.ConnectReply)
		}

		if strConnectReply.S8Error != 0 {
			return SocketError(strConnectReply.S8Error)
		}

		// NOTE: Extra data is the u16AppDataOffset member of the union in the original tstrConnectReply struct
		sockets[socket].offset = strConnectReply.U16ExtraData - protocol.M2M_HIF_HDR_OFFSET

	} else {
		err = ErrSocketDoesNotExist
	}

	return
}

func (w *WINC) Shutdown(socket Socket) (err error) {
	sockets[socket].mutex.Lock()
	defer sockets[socket].mutex.Unlock()

	if sockets[socket].inUse {
		return ErrSocketDoesNotExist
	}

	cmd := OpcodeSocketClose
	if sockets[socket].sslFlags&sslFlagsActive != 0 {
		cmd = OpcodeSocketSslClose
	}

	strClose := types.CloseCmd{
		Sock:         int8(socket),
		U16SessionID: sockets[socket].sessionId,
	}

	if err = hif.Send(GroupIP, cmd, strClose.Bytes(), nil, 0); err != nil {
		return
	}

	return
}

func (w *WINC) Secure(socket Socket) (err error) {
	sockets[socket].mutex.Lock()
	defer sockets[socket].mutex.Unlock()

	if socket < 0 || socket >= maxSocket {
		return ErrSocketInvalidArg
	} else if !sockets[socket].inUse {
		return ErrSocketInvalid
	}

	flags := sockets[socket].sslFlags
	if flags&sslFlagsDelay == 0 || flags&sslFlagsActive == 0 || sockets[socket].offset == 0 {
		return ErrSocketInvalidArg
	}

	sockets[socket].sslFlags &= sslFlagsDelay
	strConnect := types.ConnectCmd{
		Sock:         int8(socket),
		U8SslFlags:   sockets[socket].sslFlags,
		U16SessionID: sockets[socket].sessionId,
	}

	if err = hif.Send(GroupIP, OpcodeSocketSecure, strConnect.Bytes(), nil, 0); err != nil {
		return
	}

	// Wait for the response
	var strConnectReply types.ConnectReply
	select {
	case reply := <-sockets[socket].callbackChan:
		strConnectReply = reply.(types.ConnectReply)
	}

	if strConnectReply.S8Error < 0 {
		return SocketError(strConnectReply.S8Error)
	}

	// NOTE: Extra data is the u16AppDataOffset member of the union in the original tstrConnectReply struct
	sockets[socket].offset = strConnectReply.U16ExtraData - protocol.M2M_HIF_HDR_OFFSET

	return
}

func (w *WINC) Setsockopt(socket Socket, level, name int, value []byte) (err error) {
	sockets[socket].mutex.Lock()
	defer sockets[socket].mutex.Unlock()

	if socket < 0 || socket >= maxSocket {
		return ErrSocketInvalidArg
	} else if !sockets[socket].inUse {
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
				sockets[socket].sslFlags |= uint8(optVal)
			} else {
				sockets[socket].sslFlags &= ^uint8(optVal)
			}
			return
		} else if ((name == types.SO_SSL_SNI) && (len(value) < 64)) || ((name == types.SO_SSL_ALPN) && (len(value) <= 32)) {
			strSslSetSockOpt := types.SSLSetSockOptCmd{
				Sock:         int8(socket),
				U8Option:     uint8(name),
				U16SessionID: sockets[socket].sessionId,
			}

			copy(strSslSetSockOpt.Au8OptVal[:], value)

			ref, _ := strSslSetSockOpt.PassRef()
			control = C.GoBytes(unsafe.Pointer(ref), C.int(unsafe.Sizeof(*ref)))
		} else {
			return ErrSocketInvalidArg
		}
	} else if level == SolSocket && len(value) == 4 {
		strSetSockOpt := types.SetSocketOptCmd{
			U32OptionValue: protocol.ToUint32(value),
			Sock:           int8(socket),
			U8Option:       uint8(name),
			U16SessionID:   sockets[socket].sessionId,
		}

		ref, _ := strSetSockOpt.PassRef()
		control = C.GoBytes(unsafe.Pointer(ref), C.int(unsafe.Sizeof(*ref)))
	}

	if err = hif.Send(GroupIP, cmd|protocol.OpcodeReqDataPkt, control, nil, 0); err != nil {
		return ErrSocketInvalid
	}

	return
}

func (w *WINC) Send(socket Socket, buf []byte) (sz int, err error) {
	sockets[socket].mutex.Lock()
	defer sockets[socket].mutex.Unlock()

	if socket < 0 || socket >= maxSocket {
		return 0, ErrSocketInvalidArg
	} else if !sockets[socket].inUse {
		return 0, ErrSocketInvalid
	}

	cmd := OpcodeSocketSend
	offset := tcpTxPacketOffset

	if socket >= maxTcpSocket {
		offset = udpTxPacketOffset
	}

	if sockets[socket].sslFlags&uint8(sslFlagsActive) != 0 && sockets[socket].sslFlags&uint8(sslFlagsDelay) == 0 {
		cmd = OpcodeSocketSslSend
		offset = sockets[socket].offset
	}

	strSend := types.SendCmd{
		Sock:         int8(socket),
		U16DataSize:  uint16(len(buf)),
		U16SessionID: sockets[socket].sessionId,
	}

	if err = hif.Send(GroupIP, cmd|protocol.OpcodeReqDataPkt, strSend.Bytes(), buf, offset); err != nil {
		return 0, ErrSocketBufferFull
	}

	var strSendReply types.SendReply
	timer := time.NewTimer(time.Second * 30)

	// Wait for the response
	select {
	case reply := <-sockets[socket].callbackChan:
		strSendReply = reply.(types.SendReply)
	case <-timer.C:
		return 0, ErrSocketTimeout
	}

	// Check for error
	if strSendReply.S16SentBytes < 0 {
		return 0, SocketError(strSendReply.S16SentBytes)
	}

	sz = int(strSendReply.S16SentBytes)
	return
}

func (w *WINC) SendTo(socket Socket, buf []byte, addr Sockaddr) (sz int, err error) {
	sockets[socket].mutex.Lock()
	defer sockets[socket].mutex.Unlock()

	if socket < 0 || socket >= maxSocket {
		return 0, ErrSocketInvalidArg
	} else if !sockets[socket].inUse {
		return 0, ErrSocketInvalid
	}

	strSend := types.SendCmd{
		Sock:        int8(socket),
		U16DataSize: uint16(len(buf)),
		StrAddr: types.SockAddr{
			U16Family: afInet,
			U16Port:   addr.Port,
			U32IPAddr: addr.Address,
		},
		U16SessionID: sockets[socket].sessionId,
	}

	if err = hif.Send(GroupIP, OpcodeSocketSendTo|protocol.OpcodeReqDataPkt, strSend.Bytes(), buf, udpTxPacketOffset); err != nil {
		return 0, ErrSocketBufferFull
	}

	var strSendReply types.SendReply
	timer := time.NewTimer(time.Second * 30)

	// Wait for the response
	select {
	case reply := <-sockets[socket].callbackChan:
		strSendReply = reply.(types.SendReply)
	case <-timer.C:
		return 0, ErrSocketTimeout
	}

	// Check for error
	if strSendReply.S16SentBytes < 0 {
		return 0, SocketError(strSendReply.S16SentBytes)
	}

	sz = int(strSendReply.S16SentBytes)
	return
}

func (w *WINC) Recv(socket Socket, buf []byte, timeout time.Duration) (sz int, err error) {
	sockets[socket].mutex.Lock()
	defer sockets[socket].mutex.Unlock()

	if socket < 0 || socket >= maxSocket {
		return 0, ErrSocketInvalidArg
	} else if !sockets[socket].inUse {
		return 0, ErrSocketInvalid
	}

	// Only receive new bytes when the request is larger than the contents of the receiver buffer for this socket.
	if len(buf) > sockets[socket].buffer.Length() {
		cmd := OpcodeSocketRecv
		if sockets[socket].sslFlags&uint8(sslFlagsActive) != 0 && sockets[socket].sslFlags&uint8(sslFlagsDelay) == 0 {
			cmd = OpcodeSocketSslRecv
		}

		strRecv := types.RecvCmd{
			U32Timeoutmsec: uint32(timeout.Milliseconds()),
			Sock:           int8(socket),
			U16SessionID:   sockets[socket].sessionId,
			U16BufLen:      uint16(len(buf)),
		}

		if err = hif.Send(GroupIP, cmd, strRecv.Bytes(), nil, 0); err != nil {
			return 0, ErrSocketBufferFull
		}

		var strRecvReply types.RecvReply

		// Wait for the reply
		select {
		case reply := <-sockets[socket].callbackChan:
			strRecvReply = reply.(types.RecvReply)
		}

		if strRecvReply.S16RecvStatus < 0 && sockets[socket].buffer.IsEmpty() {
			return 0, SocketError(strRecvReply.S16RecvStatus)
		}
	}

	sz, _ = sockets[socket].buffer.Read(buf)
	return
}

func (w *WINC) RecvFrom(socket Socket, buf []byte, timeout time.Duration) (sz int, err error) {
	sockets[socket].mutex.Lock()
	defer sockets[socket].mutex.Unlock()

	if socket < 0 || socket >= maxSocket {
		return 0, ErrSocketInvalidArg
	} else if !sockets[socket].inUse {
		return 0, ErrSocketInvalid
	}

	// Only receive new bytes when the request is larger than the contents of the receiver buffer for this socket.
	if len(buf) > sockets[socket].buffer.Length() {
		strRecv := types.RecvCmd{
			U32Timeoutmsec: uint32(timeout.Milliseconds()),
			Sock:           int8(socket),
			U16SessionID:   sockets[socket].sessionId,
			U16BufLen:      uint16(len(buf)),
		}

		if err = hif.Send(GroupIP, OpcodeSocketRecvFrom, strRecv.Bytes(), nil, 0); err != nil {
			return
		}

		var strRecvReply types.RecvReply

		// Wait for replay
		select {
		case reply := <-sockets[socket].callbackChan:
			strRecvReply = reply.(types.RecvReply)
		}

		if strRecvReply.S16RecvStatus < 0 && sockets[socket].buffer.IsEmpty() {
			return 0, SocketError(strRecvReply.S16RecvStatus)
		}
	}

	sz, _ = sockets[socket].buffer.Read(buf)
	return
}

func (w *WINC) GetHostByName(hostname string) (address uint32, err error) {
	buf := make([]byte, hostnameMaxLength+1)
	if len(hostname) <= hostnameMaxLength {
		copy(buf, hostname)
		if err = hif.Send(GroupIP, OpcodeSocketDnsResolve, buf[:len(hostname)+1], nil, 0); err != nil {
			return
		}

		timeout := time.NewTimer(time.Second * 35)
		select {
		case strDnsReply := <-hostnameChan:
			address = strDnsReply.U32HostIP
		case <-timeout.C:
			err = ErrOperationTimeout
		}
	}

	return
}

func (w *WINC) socketCallback(id protocol.OpcodeId, sz uint16, address uint32) (data any, err error) {
	switch id {
	case OpcodeSocketAccept:
		var strAcceptReply types.AcceptReply
		if err = hif.Receive(address, strAcceptReply.Bytes(), false); err != nil {
			return
		}

		strAcceptReply.Deref()
		strAcceptReply.Free()

		sockets[strAcceptReply.SConnectedSock].offset = strAcceptReply.U16AppDataOffset
		sockets[strAcceptReply.SConnectedSock].inUse = true

		// Get unique session id
		sessionCounterMutex.Lock()
		sockets[strAcceptReply.SConnectedSock].sessionId = sessionCounter
		sessionCounter++
		sessionCounterMutex.Unlock()

		data = strAcceptReply
	case OpcodeSocketBind:
		fallthrough
	case OpcodeSocketSslBind:
		var strBindReply types.BindReply
		if err = hif.Receive(address, strBindReply.Bytes(), false); err != nil {
			return
		}

		strBindReply.Deref()
		strBindReply.Free()

		sockets[strBindReply.Sock].callbackChan <- strBindReply
		data = strBindReply
	case OpcodeSocketConnect:
		fallthrough
	case OpcodeSocketSslConnect:
		var strConnectReply types.ConnectReply
		if err = hif.Receive(address, strConnectReply.Bytes(), false); err != nil {
			return
		}

		strConnectReply.Deref()
		strConnectReply.Free()

		sockets[strConnectReply.Sock].callbackChan <- strConnectReply
		data = strConnectReply
	case OpcodeSocketListen:
		var strListenReply types.ListenReply
		if err = hif.Receive(address, strListenReply.Bytes(), false); err != nil {
			return
		}

		strListenReply.Deref()
		strListenReply.Free()

		sockets[strListenReply.Sock].callbackChan <- strListenReply
		data = strListenReply
	case OpcodeSocketRecv:
		fallthrough
	case OpcodeSocketSslRecv:
		fallthrough
	case OpcodeSocketRecvFrom:
		var strRecvReply types.RecvReply
		if err = hif.Receive(address, strRecvReply.Bytes(), false); err != nil {
			return
		}

		strRecvReply.Deref()
		strRecvReply.Free()

		if strRecvReply.Sock < 0 || strRecvReply.Sock >= maxSocket {
			return nil, ErrSocketDoesNotExist
		}

		if sockets[strRecvReply.Sock].sessionId == strRecvReply.U16SessionID {
			if strRecvReply.S16RecvStatus > 0 && strRecvReply.S16RecvStatus < int16(sz) {
				address += uint32(strRecvReply.U16DataOffset)

				buf := make([]byte, int(strRecvReply.S16RecvStatus))
				err = hif.Receive(address, buf, true)
				sockets[strRecvReply.Sock].buffer.Write(buf)
			}
		} else {
			err = hif.Receive(0, nil, true)
		}

		sockets[strRecvReply.Sock].callbackChan <- strRecvReply
		data = strRecvReply
	case OpcodeSocketSend:
		fallthrough
	case OpcodeSocketSslSend:
		fallthrough
	case OpcodeSocketSendTo:
		var strSendReply types.SendReply
		if err = hif.Receive(address, strSendReply.Bytes(), false); err != nil {
			return
		}

		strSendReply.Deref()
		strSendReply.Free()

		sockets[strSendReply.Sock].callbackChan <- strSendReply

		data = strSendReply
	case OpcodeSocketDnsResolve:
		var strDnsReply types.DnsReply
		if err = hif.Receive(address, strDnsReply.Bytes(), false); err != nil {
			return
		}

		strDnsReply.Deref()
		strDnsReply.Free()

		hostnameChan <- strDnsReply
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
