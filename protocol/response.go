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

type dmaState byte

const (
	dma_ready dmaState = 0x00
	dma_busy  dmaState = 0x01
)

type errorState byte

const (
	no_error                  errorState = 0x00
	unsupported_command       errorState = 0x01
	receiving_unexpected_data errorState = 0x02
	command_crc7_error        errorState = 0x03
	data_crc7_error           errorState = 0x04
	internal_error            errorState = 0x05
)

func (e errorState) Error() string {
	switch e {
	case unsupported_command:
		return "unsupported command"
	case receiving_unexpected_data:
		return "receiving unexpected data packet"
	case command_crc7_error:
		return "command CRC7 error"
	case data_crc7_error:
		return "data CRC7 error"
	case internal_error:
		return "internal general error"
	default:
		return "unknown error"
	}
}

type responsePacket [3]byte

func (r *responsePacket) response() byte {
	return r[0]
}

func (r *responsePacket) dma() dmaState {
	return dmaState(r[1]&0xF0) >> 4
}

func (r *responsePacket) err() error {
	e := errorState(r[1] & 0x0F)
	if e != no_error {
		return errorState(r[1] & 0x0F)
	}

	return nil
}
