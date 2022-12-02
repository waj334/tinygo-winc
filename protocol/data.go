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

// WINC default data packet size is 1024
var dataPacketSize = 1024

// SetDataPacketSize sets the data packet size returned from the device client. Call this before initialization.
func SetDataPacketSize(sz int) {
	dataPacketSize = sz
}

type dataPacket []byte

func (data dataPacket) read(t *transport, clockless, crcEnabled bool) (err error) {
	t.cs.Low()
	defer t.cs.High()
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
	t.cs.Low()
	defer t.cs.High()

	// Determine how many data packets will be sent
	count := (len(data) / dataPacketSize) + 1
	// Tell the client that this is the first packet
	seq := first

	if count == 1 {
		// This is the first and only packet to be sent
		seq = last
	}
	offset := 0

	for i := 0; i < count; i++ {
		// Transfer the header first
		if _, err := t.Transfer(byte(0xF0 | seq)); err != nil {
			return err
		}

		// Transmit a portion of slice
		if _, err := t.Write(data[offset:min(len(data), offset+dataPacketSize)]); err != nil {
			return err
		}

		if crcEnabled {
			// Send the CRC16
			// Note: Microchip's driver does not even calculate this
			t.Transfer(0)
			t.Transfer(0)
		}

		// Move offset forward
		offset += dataPacketSize

		// Indicate which packet the next is
		if i+1 != count-1 {
			// Tell the client there is more to receive
			seq = sendRecv
		} else {
			// Tell the client this the next is last packet
			seq = last
		}
	}

	return nil
}

func (data dataPacket) response(t *transport) (err error) {
	var response [3]byte

	rspLen := 3
	if t.crcEnabled {
		rspLen = 2
	}

	t.cs.Low()
	response[0], _ = t.Transfer(0)
	response[1], _ = t.Transfer(0)
	response[2], _ = t.Transfer(0)
	t.cs.High()

	if errCode := response[rspLen-1]; errCode != 0 {
		return errorState(errCode & 0x0F)
	} else if response[rspLen-2] != 0xC3 {
		return errOperationFailed
	}

	return
}
