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

	remaining := len(b)
	for remaining > 0 {
		var count int
		count, err = s.Recv(b[n:], s.recvDeadline)

		remaining -= count
		n += count

		if err == ErrSocketTimeout {
			return n, os.ErrDeadlineExceeded
		} else if err == ErrSocketConnAborted {
			return n, net.ErrClosed
		} else if err != nil {
			return n, err
		}

		// Check read deadline
		if !s.recvDeadline.IsZero() && time.Now().After(s.recvDeadline) {
			return n, os.ErrDeadlineExceeded
		}
	}

	return n, nil
}

func (s *Socket) Write(b []byte) (n int, err error) {
	// Check if the socket is valid. If not, it was likely closed
	if s.sockfd < 0 {
		return 0, net.ErrClosed
	}

	remaining := len(b)

	for remaining > 0 {
		var count int
		count, err = s.Send(b[n:], s.sendDeadline)

		remaining -= count
		n += count

		if err == ErrSocketConnAborted {
			err = net.ErrClosed
		} else if err == ErrSocketTimeout {
			err = os.ErrDeadlineExceeded
		}

		// Check write deadline
		if !s.sendDeadline.IsZero() && time.Now().After(s.sendDeadline) {
			return n, os.ErrDeadlineExceeded
		}
	}

	return n, nil
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
