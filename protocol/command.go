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

import "github.com/waj334/tinygo-winc/debug"

const (
	dmaWrite               byte = 0xC1
	dmaRead                byte = 0xC2
	registerInternalWrite  byte = 0xC3
	registerInternalRead   byte = 0xC4
	transactionTermination byte = 0xC5
	repeatDataPacket       byte = 0xC6
	dmaExtendedWrite       byte = 0xC7
	dmaExtendedRead        byte = 0xC8
	dmaSingleWordWrite     byte = 0xC9
	dmaSingleWordRead      byte = 0xCA
	softReset              byte = 0xCF

	szCommandA = 4
	szCommandB = 6
	szCommandC = 7
	szCommandD = 8
)

type commandPacket struct {
	data   [8]byte
	length uint8
}

func (cmd *commandPacket) zero() {
	for i := range cmd.data {
		cmd.data[i] = 0
	}
}

func (cmd *commandPacket) dmaSingleWordRead(address uint32) {
	cmd.length = szCommandA

	// Set the payload bytes
	cmd.data[0] = dmaSingleWordRead
	cmd.data[1] = byte(address >> 16)
	cmd.data[2] = byte(address >> 8)
	cmd.data[3] = byte(address)

}

func (cmd *commandPacket) registerInternalRead(address uint32, clockless bool) {
	cmd.length = szCommandA

	// Set the payload bytes
	cmd.data[0] = registerInternalRead
	cmd.data[1] = byte(address >> 8)
	cmd.data[2] = byte(address)
	cmd.data[3] = 0

	if clockless {
		// Set bit 15 of address to 1
		cmd.data[1] |= 1 << 7
	}
}

func (cmd *commandPacket) softReset() {
	cmd.length = szCommandA

	// Set the payload bytes
	cmd.data[0] = softReset
	cmd.data[1] = byte(0xFF)
	cmd.data[2] = byte(0xFF)
	cmd.data[3] = byte(0xFF)
}

func (cmd *commandPacket) dmaExtendedRead(address uint32, length int) {
	cmd.length = szCommandC

	cmd.data[0] = dmaExtendedRead
	cmd.data[1] = byte(address >> 16)
	cmd.data[2] = byte(address >> 8)
	cmd.data[3] = byte(address)
	cmd.data[4] = byte(length >> 16)
	cmd.data[5] = byte(length >> 8)
	cmd.data[6] = byte(length)
}

func (cmd *commandPacket) dmaExtendedWrite(address uint32, length int) {
	cmd.length = szCommandC

	cmd.data[0] = dmaExtendedWrite
	cmd.data[1] = byte(address >> 16)
	cmd.data[2] = byte(address >> 8)
	cmd.data[3] = byte(address)
	cmd.data[4] = byte(length >> 16)
	cmd.data[5] = byte(length >> 8)
	cmd.data[6] = byte(length)
}

func (cmd *commandPacket) registerInternalWrite(address, value uint32, clockless bool) {
	cmd.length = szCommandC

	// Set the payload bytes
	cmd.data[0] = registerInternalWrite
	cmd.data[1] = byte(address >> 8)
	cmd.data[2] = byte(address)
	cmd.data[3] = byte(value >> 24)
	cmd.data[4] = byte(value >> 16)
	cmd.data[5] = byte(value >> 8)
	cmd.data[6] = byte(value)

	if clockless {
		// Set bit 15 of address to 1
		cmd.data[1] |= 1 << 7
	}
}

func (cmd *commandPacket) dmaSingleWordWrite(address, value uint32) {
	cmd.length = szCommandD

	// Set the payload bytes
	cmd.data[0] = dmaSingleWordWrite
	cmd.data[1] = byte(address >> 16)
	cmd.data[2] = byte(address >> 8)
	cmd.data[3] = byte(address)
	cmd.data[4] = byte(value >> 24)
	cmd.data[5] = byte(value >> 16)
	cmd.data[6] = byte(value >> 8)
	cmd.data[7] = byte(value)
}

func (cmd *commandPacket) calculateCRC8() byte {
	return crc7(cmd.data[:cmd.length]) << 1
}

func (cmd *commandPacket) write(t *transport) (err error) {
	t.chipSelect(false)
	defer t.chipSelect(true)

	debug.DEBUG("PROTOCOL: Sending command %X: LEN=%v %v", cmd.data[0], cmd.length, cmd.data)
	defer debug.DEBUG("PROTOCOL: Done sending command %X", cmd.data[0])

	// Write the data payload
	if _, err = t.Write(cmd.data[:cmd.length]); err != nil {
		return
	}

	// Write the CRC byte if enabled
	if t.crcEnabled {
		if _, err = t.Transfer(cmd.calculateCRC8()); err != nil {
			return err
		}
	}

	return
}

func (cmd *commandPacket) response(t *transport, clockless bool) (err error) {
	t.chipSelect(false)
	defer t.chipSelect(true)

	debug.DEBUG("PROTOCOL: Receiving command %X response", cmd.data[0])
	defer debug.DEBUG("PROTOCOL: Done receiving command %X response", cmd.data[0])

	if cmd.data[0] == softReset || cmd.data[0] == transactionTermination || cmd.data[0] == repeatDataPacket {
		// Attempt to read and return any error
		// NOTE: Nothing was done with the data read below. This read the additional leading byte I observed when
		//       testing the reset cycle. Fixed the issue when resetting.
		if _, err = t.Transfer(0); err != nil {
			return
		}
	}

	var response byte
	response, err = t.Transfer(0)
	debug.DEBUG("PROTOCOL: response: %X", response)

	if err != nil {
		return err
	} else if response != cmd.data[0] {
		return errOperationFailed
	}

	var code byte
	code, err = t.Transfer(0)
	debug.DEBUG("PROTOCOL: code: %X", code)

	if err != nil {
		return err
	} else if code != 0 {
		return errorState(code & 0x0F)
	}

	return
}
