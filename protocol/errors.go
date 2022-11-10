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

import "errors"

// TODO: Reduce these to a few generic errors to size program memory
var (
	errExceededMaxRetries = errors.New("exceeded max retries")
	errProtocolFailed     = errors.New("failed internal read protocol")
	errInvalidPacketSize  = errors.New("invalid packet size")
	errFailedDataWrite    = errors.New("failed data write response")
	errNoStartHeader      = errors.New("never got start header")
	errMessageTooLong     = errors.New("message is too long")
	errBadMemoryAlloc     = errors.New("memory allocation error")
	errOperationFailed    = errors.New("operation failed")
	errUnknownCallback    = errors.New("unknown callback")

	errIncompatibleVersion    = errors.New("device reported incompatible version")
	errBootromFailed          = errors.New("bootrom failed")
	errFirmwareFailed         = errors.New("firmware failed to start")
	errFirmwareLoadFailed     = errors.New("failed to load firmware from flash")
	errFirmwareTimeout        = errors.New("timeout while waiting for firmware start")
	errInterruptsEnableFailed = errors.New("failed to enable interrupts on the device")
	errChipWakeFail           = errors.New("failed to wake the chip")
)
