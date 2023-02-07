package winc

import (
	"context"
	"math"
	"net"
	"net/url"
	"unsafe"
)

func (w *WINC) Dial(network, address string) (conn net.Conn, err error) {
	return w.DialContext(context.Background(), network, address, false)
}

func (w *WINC) DialTLS(network, address string) (conn net.Conn, err error) {
	return w.DialContext(context.Background(), network, address, true)
}

func (w *WINC) DialContext(ctx context.Context, network, address string, tls bool) (conn net.Conn, err error) {
	var socket *Socket
	var addr net.Addr
	var ip uint32
	var uri *url.URL

	// Parse the url
	if uri, err = url.Parse(network + "://" + address); err != nil {
		return nil, err
	}

	// Parse the port number string
	port := uint16(Atoi(uri.Port()))

	// Perform DNS lookup
	if ip, err = w.GetHostByName(uri.Hostname()); err != nil {
		return nil, err
	}

	config := SocketConfigSslOff
	if tls {
		config = SocketConfigSslOn
	}

	// Create the respective socket type
	if network == "tcp" {
		if socket, err = w.Socket(SocketTypeStream, config); err != nil {
			return nil, err
		}

		addr = &TCPAddr{
			Family:    afInet,
			Port:      port,
			IPAddress: ip,
		}
	} else if network == "udp" {
		if socket, err = w.Socket(SocketTypeDatagram, config); err != nil {
			return nil, err
		}

		addr = &UDPAddr{
			Family:    afInet,
			Port:      port,
			IPAddress: ip,
		}
	} else {
		return nil, &net.AddrError{
			Err:  "unsupported network scheme",
			Addr: uri.String(),
		}
	}

	if tls {
		val := 1
		if err = socket.Setsockopt(SslSocketLevel, SslEnableSessionCaching, unsafe.Slice((*uint8)(unsafe.Pointer(&val)), 4)); err != nil {
			panic(err)
		}

		if err = socket.Setsockopt(SslSocketLevel, SslBypassX509Verification, unsafe.Slice((*uint8)(unsafe.Pointer(&val)), 4)); err != nil {
			panic(err)
		}
	}

	// Connect to the url
	if err = socket.Connect(addr); err != nil {
		// Free the socket
		socket.Shutdown()

		return nil, err
	}

	return socket, nil
}

func Atoi(input string) (result int) {
	digit := 0
	for i := len(input) - 1; i >= 0; i-- {
		result += int(input[i]-48) * int(math.Pow10(digit))
		digit++
	}

	return
}
