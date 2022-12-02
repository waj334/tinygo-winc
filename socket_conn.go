package winc

import (
	"net"
	"os"
	"time"
)

func (s *Socket) Read(b []byte) (n int, err error) {
	// Check if the socket is valid. If not, it was likely closed
	if s.sockfd < 0 {
		return 0, net.ErrClosed
	}

	n, err = s.Recv(b, s.recvDeadline)
	if err == ErrSocketTimeout {
		err = os.ErrDeadlineExceeded
	} else if err == ErrSocketConnAborted {
		err = net.ErrClosed
	}

	return
}

func (s *Socket) Write(b []byte) (n int, err error) {
	// Check if the socket is valid. If not, it was likely closed
	if s.sockfd < 0 {
		return 0, net.ErrClosed
	}

	if n, err = s.Send(b, time.Now()); err == ErrSocketConnAborted {
		err = net.ErrClosed
	} else if err == ErrSocketTimeout {
		err = os.ErrDeadlineExceeded
	}

	return
}

func (s *Socket) Close() error {
	// Check if the socket is valid. If not, it was likely closed
	if s.sockfd < 0 {
		return net.ErrClosed
	}

	return s.Shutdown()
}

func (s *Socket) LocalAddr() net.Addr {
	return &s.driver.ipAddr
}

func (s *Socket) RemoteAddr() net.Addr {
	return s.addr
}

func (s *Socket) SetDeadline(t time.Time) error {
	// Check if the socket is valid. If not, it was likely closed
	if s.sockfd < 0 {
		return net.ErrClosed
	}

	s.recvDeadline = t
	return nil
}

func (s *Socket) SetReadDeadline(t time.Time) error {
	// Check if the socket is valid. If not, it was likely closed
	if s.sockfd < 0 {
		return net.ErrClosed
	}

	s.recvDeadline = t
	return nil
}

func (s *Socket) SetWriteDeadline(t time.Time) error {
	// Check if the socket is valid. If not, it was likely closed
	if s.sockfd < 0 {
		return net.ErrClosed
	}

	s.sendDeadline = t
	return nil
}
