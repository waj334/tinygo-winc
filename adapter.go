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
	"math"
	"time"

	"github.com/waj334/tinygo-winc/protocol"
)

var currentSocket Socket = SocketInvalid

func (w *WINC) ConnectToAccessPoint(ssid, pass string, timeout time.Duration) (err error) {
	return w.WifiConnectPsk(WifiConnectionSettings{
		Ssid:       ssid,
		Channel:    WifiChannelAll,
		Passphrase: pass,
		Storage:    WifiCredSaveEncrypted,
		security:   WifiSecurityWpaPsk,
	})
}

func (w *WINC) Disconnect() (err error) {
	return w.WifiDisconnect()
}

func (w *WINC) GetClientIP() (val string, err error) {
	val += protocol.Uitoa(uint(w.ipAddress&0xFF)) + "."
	val += protocol.Uitoa(uint((w.ipAddress>>8)&0xFF)) + "."
	val += protocol.Uitoa(uint((w.ipAddress>>16)&0xFF)) + "."
	val += protocol.Uitoa(uint((w.ipAddress >> 24) & 0xFF))

	return
}

func (w *WINC) GetDNS(domain string) (val string, err error) {
	var address uint32
	if address, err = w.GetHostByName(domain); err != nil {
		return
	}

	val += protocol.Uitoa(uint(address&0xFF)) + "."
	val += protocol.Uitoa(uint((address>>8)&0xFF)) + "."
	val += protocol.Uitoa(uint((address>>16)&0xFF)) + "."
	val += protocol.Uitoa(uint((address >> 24) & 0xFF))

	return
}

func (w *WINC) ConnectTCPSocket(addr, port string) (err error) {
	if currentSocket != SocketInvalid {
		// Disconnect the current socket
		if err = w.DisconnectSocket(); err != nil {
			return
		}
	}

	// Create a new socket
	if currentSocket, err = w.Socket(SocketTypeStream, SocketConfigSslOff); err != nil {
		return
	}

	// Attempt to connect to the address
	uAddr := ParseAddress(addr)
	uPort := Htons(uint16(Atoi(port)))

	if err = w.Connect(currentSocket, Sockaddr{
		Port:    uPort,
		Address: uAddr,
	}); err != nil {
		return
	}

	return
}
func (w *WINC) ConnectSSLSocket(addr, port string) (err error) {
	if currentSocket != SocketInvalid {
		// Disconnect the current socket
		if err = w.DisconnectSocket(); err != nil {
			return
		}
	}

	// Create a new socket
	if currentSocket, err = w.Socket(SocketTypeStream, SocketConfigSslOn); err != nil {
		return
	}

	// Attempt to connect to the address
	uAddr := ParseAddress(addr)
	uPort := Htons(uint16(Atoi(port)))

	if err = w.Connect(currentSocket, Sockaddr{
		Port:    uPort,
		Address: uAddr,
	}); err != nil {
		return
	}

	return
}
func (w *WINC) ConnectUDPSocket(addr, sendport, listenport string) (err error) {
	if currentSocket != SocketInvalid {
		// Disconnect the current socket
		if err = w.DisconnectSocket(); err != nil {
			return
		}
	}

	// Create a new socket
	if currentSocket, err = w.Socket(SocketTypeDatagram, SocketConfigSslOff); err != nil {
		return
	}

	// TODO: Store send and listen ports
	/*
		// Attempt to connect to the address
		uAddr := ParseAddress(addr)
		uSendPort := Htons(uint16(Atoi(sendport)))
		uListenPort := Htons(uint16(Atoi(listenport)))

		if err = w.Connect(currentSocket, Sockaddr{
			Port:    uSendPort,
			Address: uAddr,
		}); err != nil {
			return
		}
	*/

	return
}
func (w *WINC) DisconnectSocket() (err error) {
	err = w.Shutdown(currentSocket)
	currentSocket = SocketInvalid
	return
}

func (w *WINC) StartSocketSend(size int) (err error) {
	return
}

func (w *WINC) Write(b []byte) (n int, err error) {
	return w.Send(currentSocket, b)
}

func (w *WINC) ReadSocket(b []byte) (n int, err error) {
	if n, err = w.Recv(currentSocket, b, time.Second); err == ErrSocketTimeout {
		// Ignore this error
		err = nil
	}
	return
}

func (w *WINC) IsSocketDataAvailable() bool {
	return true
	//return !sockets[currentSocket].buffer.IsEmpty()
}

func (w *WINC) Response(timeout int) (data []byte, err error) {
	return
}

func ParseAddress(addr string) (result uint32) {
	digit := 0
	value := uint32(0)
	offset := 24
	for i := len(addr) - 1; i >= 0; i-- {
		if addr[i] == '.' {
			result |= value << offset
			offset -= 8
			value = 0
			digit = 0
		} else {
			number := uint32(addr[i] - 48)
			if number < 0 || number > 9 {
				return 0
			}

			value += number * uint32(math.Pow10(digit))
			if value > 255 {
				return 0
			}

			digit++
		}
	}

	result |= value
	return
}

func Atoi(input string) (result int) {
	digit := 0
	for i := len(input) - 1; i >= 0; i-- {
		result += int(input[i]-48) * int(math.Pow10(digit))
		digit++
	}

	return
}
