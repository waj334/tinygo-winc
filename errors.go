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
	"errors"
)

var (
	ErrInvalidParameter   = errors.New("invalid parameter")
	ErrOperationTimeout   = errors.New("operation timed out")
	ErrNoAvailableSocket  = errors.New("no available socket")
	ErrSocketDoesNotExist = errors.New("socket does not exist")
	ErrUnknown            = errors.New("unknown error occurred")

	ErrSocketInvalidAddress     = SocketError(-1)
	ErrSocketAddrAlreadyInUse   = SocketError(-2)
	ErrSocketMaxTcpSock         = SocketError(-3)
	ErrSocketMaxUdpSock         = SocketError(-4)
	ErrSocketInvalidArg         = SocketError(-6)
	ErrSocketMaxListenSock      = SocketError(-7)
	ErrSocketInvalid            = SocketError(-9)
	ErrSocketAddrIsRequired     = SocketError(-11)
	ErrSocketConnAborted        = SocketError(-12)
	ErrSocketTimeout            = SocketError(-13)
	ErrSocketBufferFull         = SocketError(-14)
	ErrSocketFuncNotImplemented = SocketError(-99)
)

type SocketError int8

func (s SocketError) Error() (err string) {
	switch s {
	case -1:
		err = "socket address is invalid"
	case -2:
		err = "socket operation cannot bind on the given address"
	case -3:
		err = "exceeded the maximum number of TCP sockets"
	case -4:
		err = "exceeded the maximum number of UDP sockets"
	case -6:
		err = "an invalid option passed to a socket function"
	case -7:
		err = "exceeded the maximum number of passive TCP listening sockets"
	case -9:
		err = "the request socket operation is not valid in the current socket state"
	case -11:
		err = "destination address is required"
	case -12:
		err = "socket is closed (reset) by the peer"
	case -13:
		err = "the pending socket operation has timed out"
	case -14:
		err = "no buffer space available for the requested socket operation"
	}

	return
}
