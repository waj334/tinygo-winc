package winc

import (
	"context"
	"math"
	"net"
	"net/url"
)

func (w *WINC) Dial(network, address string) (conn net.Conn, err error) {
	return w.DialContext(context.Background(), network, address)
}

func (w *WINC) DialContext(ctx context.Context, network, address string) (conn net.Conn, err error) {
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

	// Create the respective socket type
	if network == "tcp" {
		if socket, err = w.Socket(SocketTypeStream, SocketConfigSslOff); err != nil {
			return nil, err
		}

		addr = &TCPAddr{
			U16Family: afInet,
			U16Port:   port,
			U32IPAddr: ip,
		}
	} else if network == "udp" {
		if socket, err = w.Socket(SocketTypeDatagram, SocketConfigSslOff); err != nil {
			return nil, err
		}

		addr = &UDPAddr{
			U16Family: afInet,
			U16Port:   port,
			U32IPAddr: ip,
		}
	} else {
		return nil, &net.AddrError{
			Err:  "unsupported network scheme",
			Addr: uri.String(),
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
