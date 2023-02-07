package winc

import (
	"github.com/waj334/tinygo-winc/protocol"
)

type TCPAddr SocketAddress

func (a *TCPAddr) Network() string {
	return "tcp"
}

func (a *TCPAddr) String() (address string) {
	// Format the IP as a string
	address += protocol.Uitoa(uint(a.IPAddress&0xFF)) + "."
	address += protocol.Uitoa(uint((a.IPAddress>>8)&0xFF)) + "."
	address += protocol.Uitoa(uint((a.IPAddress>>16)&0xFF)) + "."
	address += protocol.Uitoa(uint((a.IPAddress >> 24) & 0xFF))

	// Append the port to the string
	address += ":" + protocol.Itoa(int(a.Port))

	return
}

type UDPAddr SocketAddress

func (a *UDPAddr) Network() string {
	return "udp"
}

func (a *UDPAddr) String() (address string) {
	// Format the IP as a string
	address += protocol.Uitoa(uint(a.IPAddress&0xFF)) + "."
	address += protocol.Uitoa(uint((a.IPAddress>>8)&0xFF)) + "."
	address += protocol.Uitoa(uint((a.IPAddress>>16)&0xFF)) + "."
	address += protocol.Uitoa(uint((a.IPAddress >> 24) & 0xFF))

	// Append the port to the string
	address += ":" + protocol.Itoa(int(a.Port))

	return
}
