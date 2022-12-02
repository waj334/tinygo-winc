package winc

import (
	"net"
)

func (s *Socket) Accept() (net.Conn, error) {
	// Check if the socket is valid. If not, it was likely closed
	if s.sockfd < 0 {
		return nil, net.ErrClosed
	}

	// Wait for socket to be ready
	connectedSockfd := <-s.acceptChan

	// Return the error if the returned socket is < 0
	if connectedSockfd < 0 {
		return nil, &net.OpError{
			Op:     "Accept",
			Net:    s.addr.Network(),
			Source: s.addr,
			Addr:   &s.driver.ipAddr,
			Err:    SocketError(connectedSockfd),
		}
	}

	// Return the connected socket
	return s.driver.sockets[connectedSockfd], nil
}

func (s *Socket) Addr() net.Addr {
	return &s.driver.ipAddr
}
