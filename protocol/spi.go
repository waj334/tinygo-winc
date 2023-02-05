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

package protocol

import (
	"github.com/waj334/tinygo-winc/debug"
	"sync"
	"time"

	"machine"

	"tinygo.org/x/drivers"
)

type transport struct {
	spi      drivers.SPI
	cs       machine.Pin
	busMutex sync.Mutex
	spiMutex sync.Mutex

	crcEnabled bool
}

func (t *transport) init() (err error) {
	t.crcEnabled = true

	var result uint32
	result, err = t.ReadRegister(_NMI_SPI_PROTOCOL_CONFIG)
	if err != nil {
		// Try again with CRC disabled
		t.crcEnabled = false
		result, err = t.ReadRegister(_NMI_SPI_PROTOCOL_CONFIG)
		if err != nil {
			return errProtocolFailed
		}
	}

	// NOTE: I think CRC might be required for some start up sequence. The block below disables it.
	if t.crcEnabled {
		// disable crc
		result &= 0x2
		result &= 0x8F
		result |= 0x5 << 4

		err = t.WriteRegister(_NMI_SPI_PROTOCOL_CONFIG, result)
		if err != nil {
			return err
		}

		t.crcEnabled = false
	}

	if err = t.initPacketSize(); err != nil {
		return err
	}

	return nil
}

func (t *transport) Transfer(b byte) (byte, error) {
	t.spiMutex.Lock()
	defer t.spiMutex.Unlock()

	return t.spi.Transfer(b)
}

func (t *transport) Write(b []byte) (n int, err error) {
	t.spiMutex.Lock()
	defer t.spiMutex.Unlock()

	err = t.spi.Tx(b, nil)
	n = len(b)

	return
}

func (t *transport) Read(b []byte) (n int, err error) {
	t.spiMutex.Lock()
	defer t.spiMutex.Unlock()

	err = t.spi.Tx(nil, b)
	n = len(b)

	return
}

func (t *transport) initPacketSize() error {
	// Set the packet size
	result, err := t.ReadRegister(_SPI_BASE + 0x24)
	if err != nil {
		return err
	}

	result &= ^uint32(0x7 << 4)
	switch dataPacketSize {
	case 256:
		result |= 0 << 4
	case 512:
		result |= 1 << 4
	case 1024:
		result |= 2 << 4
	case 2048:
		result |= 3 << 4
	case 4096:
		result |= 4 << 4
	case 8192:
		result |= 5 << 5
	default:
		return errInvalidPacketSize
	}

	// Write the packet size setting
	err = t.WriteRegister(_SPI_BASE+0x24, result)
	if err != nil {
		return err
	}

	return nil
}

func (t *transport) ReadRegister(address uint32) (result uint32, err error) {
	t.busMutex.Lock()
	defer t.busMutex.Unlock()

	for retry := 0; retry < 10; retry++ {
		cmd := commandPacket{}
		clockless := false

		if address <= 0xFF {
			clockless = true
			cmd.registerInternalRead(address, clockless)
		} else {
			cmd.dmaSingleWordRead(address)
		}

		// Send command
		if err = cmd.write(t); err != nil {
			time.Sleep(time.Millisecond)
			if err = t.internalReset(); err != nil {
				return 0, err
			}
			time.Sleep(time.Millisecond)
			continue
		}

		// Wait for command response
		if err = cmd.response(t, clockless); err != nil {
			time.Sleep(time.Millisecond)
			if err = t.internalReset(); err != nil {
				return 0, err
			}
			time.Sleep(time.Millisecond)
			continue
		}

		// Receive the data response
		var buf [4]byte
		var data dataPacket = buf[:]
		if err = data.read(t, clockless, t.crcEnabled); err != nil {
			time.Sleep(time.Millisecond)
			if err = t.internalReset(); err != nil {
				return 0, err
			}
			time.Sleep(time.Millisecond)
			continue
		}

		// Set each byte of the integer result
		result = uint32(data[0]) |
			uint32(data[1])<<8 |
			uint32(data[2])<<16 |
			uint32(data[3])<<24

		break
	}

	return
}

func (t *transport) WriteRegister(address, value uint32) (err error) {
	t.busMutex.Lock()
	defer t.busMutex.Unlock()

	for retry := 0; retry < 10; retry++ {
		cmd := commandPacket{}
		clockless := false
		if address <= 0x30 {
			// Clockless write
			clockless = true
			cmd.registerInternalWrite(address, value, clockless)
		} else {
			cmd.dmaSingleWordWrite(address, value)
		}

		// Send command
		if err = cmd.write(t); err != nil {
			time.Sleep(time.Millisecond)
			if err = t.internalReset(); err != nil {
				return err
			}
			time.Sleep(time.Millisecond)
			continue
		}

		// Stop if sending reset command
		if address == _rNMI_GLB_RESET {
			return nil
		}

		// Wait for command response
		if err = cmd.response(t, clockless); err != nil {
			time.Sleep(time.Millisecond)
			if err = t.internalReset(); err != nil {
				return err
			}
			time.Sleep(time.Millisecond)
			continue
		}

		break
	}

	return
}

func (t *transport) ReadBlock(address uint32, data []byte) (err error) {
	t.busMutex.Lock()
	defer t.busMutex.Unlock()

	buf := data

	// The minimum block size is 2 bytes
	if len(data) == 1 {
		buf = make([]byte, 2)
	}

	for retry := 0; retry < 10; retry++ {
		cmd := commandPacket{}
		var pkt dataPacket = buf

		// Format the DMA extended read command
		cmd.dmaExtendedRead(address, len(buf))

		// Write the command
		if err = cmd.write(t); err != nil {
			if err = t.internalReset(); err != nil {
				return err
			}
			continue
		}

		// Get response to command
		if err = cmd.response(t, false); err != nil {
			time.Sleep(time.Millisecond)
			if err = t.internalReset(); err != nil {
				return err
			}
			time.Sleep(time.Millisecond)
			continue
		}

		// Receive the data
		if err = pkt.read(t, false, t.crcEnabled); err != nil {
			time.Sleep(time.Millisecond)
			if err = t.internalReset(); err != nil {
				return err
			}
			time.Sleep(time.Millisecond)
			continue
		}

		// Copy into input buffer
		// NOTE: This accounts for when the input buffer size is 1 byte
		copy(data, buf)

		break
	}

	return
}

func (t *transport) WriteBlock(address uint32, data []byte) (err error) {
	debug.DEBUG("Transport: WriteBlock - BEGIN")
	defer debug.DEBUG("Transport: WriteBlock - END")

	t.busMutex.Lock()
	defer t.busMutex.Unlock()

	for len(data) > 0 {
		chunk := data
		if len(chunk) > 2040 {
			chunk = chunk[:2040]
			data = data[2040:]
		} else {
			data = []byte{}
		}

		if err = t.writeBlockInternal(address, chunk); err != nil {
			return
		}

		// Advance address
		address += uint32(len(chunk))
	}

	return
}

func (t *transport) reset() (err error) {
	t.busMutex.Lock()
	defer t.busMutex.Unlock()

	return t.internalReset()
}

func (t *transport) writeBlockInternal(address uint32, data []byte) (err error) {
	// The minimum block size is 2 bytes
	buf := data
	if len(data) == 1 {
		buf = make([]byte, 2)
		buf[0] = data[0]
	}

	cmd := commandPacket{}
	cmd.dmaExtendedWrite(address, len(buf))

	pkt := dataPacket(data)

	for retry := 0; retry < 10; retry++ {
		debug.DEBUG("Transport: Attempt %v - BEGIN", retry)

		// Send the command
		if err = cmd.write(t); err != nil {
			goto reset
		}

		// Wait for the response
		if err = cmd.response(t, false); err != nil {
			goto reset
		}

		// Write the data
		if err = pkt.write(t, t.crcEnabled); err != nil {
			goto reset
		}

		// read the data response
		if err = pkt.response(t); err != nil {
			goto reset
		}

		// Stop the loop if there was no failed attempt
		return
	reset:
		time.Sleep(time.Millisecond)
		t.internalReset()
		time.Sleep(time.Millisecond)
		debug.DEBUG("Transport: Attempt %v - END", retry)
	}

	return
}

func (t *transport) internalReset() (err error) {
	debug.DEBUG("PROTOCOL: internalReset - BEGIN")
	defer debug.DEBUG("PROTOCOL: internalReset - END")

	// NOTE: Do not lock mutex in this function
	cmd := commandPacket{}
	cmd.softReset()

	// Write the command
	if err = cmd.write(t); err != nil {
		return
	}

	if err = cmd.response(t, false); err != nil {
		return
	}

	time.Sleep(time.Millisecond * 100)

	return
}

func (t *transport) chipSelect(enable bool) {
	if t.cs != machine.NoPin && t.cs != 0 {
		if enable {
			t.cs.High()
		} else {
			t.cs.Low()
		}
	}
}
