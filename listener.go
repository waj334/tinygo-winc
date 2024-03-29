package winc

import (
	"net"
	"net/url"
)

func (w *WINC) Listen(network, address string) (listener net.Listener, err error) {
	var socket *Socket
	var addr net.Addr
	var uri *url.URL

	// Parse the url
	if uri, err = url.Parse(network + "://" + address); err != nil {
		return nil, err
	}

	// Parse the port number string
	port := Atoi(uri.Port())

	// Create the respective socket type
	if network == "tcp" {
		if socket, err = w.Socket(SocketTypeStream, SocketConfigSslOff); err != nil {
			return nil, err
		}

		addr = &net.TCPAddr{
			IP:   []byte{0, 0, 0, 0},
			Port: port,
		}
	} else if network == "udp" {
		if socket, err = w.Socket(SocketTypeDatagram, SocketConfigSslOff); err != nil {
			return nil, err
		}

		addr = &net.UDPAddr{
			IP:   []byte{0, 0, 0, 0},
			Port: port,
		}
	} else {
		return nil, &net.AddrError{
			Err:  "unsupported network scheme",
			Addr: uri.String(),
		}
	}

	// Bind the socket to the listen address
	if err = socket.Bind(addr); err != nil {
		return nil, err
	}

	// Begin listening to the socket
	if err = socket.Listen(1); err != nil {
		return nil, err
	}

	return socket, nil
}
