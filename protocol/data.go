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
)

// WINC default data packet size is 1024
var dataPacketSize = 1024

// SetDataPacketSize sets the data packet size returned from the device client. Call this before initialization.
func SetDataPacketSize(sz int) {
	dataPacketSize = sz
}

type dataPacket []byte

func (data dataPacket) read(t *transport, clockless, crcEnabled bool) (err error) {
	t.chipSelect(false)
	defer t.chipSelect(true)
	// Determine how many data packets will be read
	count := (len(data) / dataPacketSize) + 1

	offset := 0
	for i := 0; i < count; i++ {
		// Receive header
		for ii := 0; ii < 10; ii++ {
			var header byte
			if header, err = t.Transfer(0); err != nil {
				return err
			}

			if header&0xF0 != 0xF0 {
				continue
			} else {
				goto readData
			}
		}

		return errNoStartHeader

	readData:
		// Receive the data
		if _, err = t.Read(data[offset:min(len(data), offset+dataPacketSize)]); err != nil {
			return
		}

		// CRC is only available during clocked reads
		if !clockless && crcEnabled {
			crc := make([]byte, 2, 2)
			crc[0], _ = t.Transfer(0)
			crc[1], _ = t.Transfer(0)
		}

		// Advance the offset
		offset += dataPacketSize
	}

	return nil
}

func (data dataPacket) write(t *transport, crcEnabled bool) error {
	t.chipSelect(false)
	defer t.chipSelect(true)

	debug.DEBUG("PROTOCOL: Writing data")
	defer debug.DEBUG("PROTOCOL: Done writing data")

	// Determine how many data packets will be sent
	count := (len(data) / dataPacketSize) + 1

	debug.DEBUG("PROTOCOL: Will write %v data packets", count)

	chunk := []byte(data)

	for i := 0; i < count; i++ {
		var seq sequence

		// Which packet is this?
		if i == 0 {
			seq = first
		} else if i == count-1 {
			seq = last
		} else {
			seq = sendRecv
		}

		// Transfer the header first
		if _, err := t.Transfer(byte(0xF0 | seq)); err != nil {
			return err
		}

		buf := chunk
		if len(chunk) > dataPacketSize {
			// Limit to the maximum data packet size
			buf = buf[:dataPacketSize]

			// Advance to the next chunk
			chunk = chunk[dataPacketSize:]
		}

		// Transmit a slice of the data
		if _, err := t.Write(buf); err != nil {
			return err
		}

		if crcEnabled {
			// Send the CRC16
			// Note: Microchip's driver does not even calculate this
			t.Transfer(0)
			t.Transfer(0)
		}

		debug.DEBUG("PROTOCOL: Wrote data packet %v", i)
	}

	return nil
}

func (data dataPacket) response(t *transport) (err error) {
	var response [3]byte

	debug.DEBUG("PROTOCOL: Receiving data response")
	defer debug.DEBUG("PROTOCOL: Done receiving data response")

	rspLen := 3
	if t.crcEnabled {
		rspLen = 2
	}

	t.chipSelect(false)
	response[0], _ = t.Transfer(0)
	response[1], _ = t.Transfer(0)
	response[2], _ = t.Transfer(0)
	t.chipSelect(true)

	debug.DEBUG("PROTOCOL: Data response packet %v", response)

	if errCode := response[rspLen-1]; errCode != 0 {
		debug.DEBUG("PROTOCOL: Data error %X", errCode)
		return errorState(errCode & 0x0F)
	} else if response[rspLen-2] != 0xC3 {
		debug.DEBUG("PROTOCOL: Data response code %X", response[rspLen-2])
		return errOperationFailed
	}

	return
}
