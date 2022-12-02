package winc

import (
	"github.com/waj334/tinygo-winc/protocol"
	"github.com/waj334/tinygo-winc/protocol/types"
)

type TCPAddr types.SockAddr

func (a *TCPAddr) Network() string {
	return "tcp"
}

func (a *TCPAddr) String() (address string) {
	// Format the IP as a string
	address += protocol.Uitoa(uint(a.U32IPAddr&0xFF)) + "."
	address += protocol.Uitoa(uint((a.U32IPAddr>>8)&0xFF)) + "."
	address += protocol.Uitoa(uint((a.U32IPAddr>>16)&0xFF)) + "."
	address += protocol.Uitoa(uint((a.U32IPAddr >> 24) & 0xFF))

	// Append the port to the string
	address += ":" + protocol.Itoa(int(a.U16Port))

	return
}

type UDPAddr types.SockAddr

func (a *UDPAddr) Network() string {
	return "udp"
}

func (a *UDPAddr) String() (address string) {
	// Format the IP as a string
	address += protocol.Uitoa(uint(a.U32IPAddr&0xFF)) + "."
	address += protocol.Uitoa(uint((a.U32IPAddr>>8)&0xFF)) + "."
	address += protocol.Uitoa(uint((a.U32IPAddr>>16)&0xFF)) + "."
	address += protocol.Uitoa(uint((a.U32IPAddr >> 24) & 0xFF))

	// Append the port to the string
	address += ":" + protocol.Itoa(int(a.U16Port))

	return
}
