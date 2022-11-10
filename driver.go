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
	"sync"
	"time"

	"github.com/waj334/tinygo-winc/protocol"
	"github.com/waj334/tinygo-winc/protocol/hal"
)

var (
	initialized       = false
	hif               protocol.Hif
	isrSignal         chan bool
	isrShutdownSignal chan bool
)

type WINC struct {
	SPI       hal.SPI
	CS        hal.Pin
	EnablePin hal.Pin
	ResetPin  hal.Pin

	SetIRQFunc   hal.SetInterruptHandlerFunc
	ResetIRQFunc hal.ResetInterruptHandlerFunc

	wifiState WifiState
	ipAddress uint32

	mutex sync.Mutex
}

func (w *WINC) Initialize() (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !initialized {
		// Create the hardware interface abstraction layer
		hif = protocol.CreateHif(w.SPI, w.CS)

		// Initialize the HAL
		err = hif.Init()
		if err != nil {
			return
		}

		// Register interrupt callbacks
		hif.RegisterCallback(GroupWIFI, w.wifiCallback)
		hif.RegisterCallback(GroupIP, w.socketCallback)

		// Start the interrupt service (go)routine
		isrSignal = make(chan bool, 1)
		isrShutdownSignal = make(chan bool, 1)
		go w.isr(isrSignal)

		// set up the interrupt
		if err = w.SetIRQFunc(); err != nil {
			return
		}

		initialized = true
	}

	return
}

func (w *WINC) OpenEventChannel() <-chan protocol.Event {
	return hif.OpenEventChannel()
}

// Reset the SoC by driving chip enable and reset pins low then high
func (w *WINC) Reset() {
	println("driver: About to reset WINC driver...")
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Unset the interrupt handler
	w.ResetIRQFunc()

	// Stop the ISR routine
	if isrSignal != nil {
		select {
		case isrSignal <- false:
			close(isrSignal)
		default:
			close(isrSignal)
		}

		// Wait for isr to be fully shutdown
		<-isrShutdownSignal
		close(isrShutdownSignal)

		isrSignal = nil
	}

	println("driver: Closed ISR signal channel")

	// Shutdown Hif
	hif.Shutdown()

	println("driver: Shutdown HIF")

	// Drive the pins low
	w.EnablePin.Low()
	w.ResetPin.Low()

	// Sleep
	time.Sleep(time.Millisecond * 100)

	// Re-enable the pins
	w.EnablePin.High()
	time.Sleep(time.Millisecond * 10)
	w.ResetPin.High()
	time.Sleep(time.Millisecond * 10)

	initialized = false

	println("driver: WINC driver reset complete")
}

func (w *WINC) SetGPIODirection(gpio GPIOType, direction GPIODirection) error {
	return hif.SetGPIODirection(uint8(gpio), uint8(direction))
}

func (w *WINC) SetGPIOState(gpio GPIOType, state GPIOState) error {
	return hif.SetGPIOValue(uint8(gpio), uint8(state))
}

func (w *WINC) GetGPIOState(gpio GPIOType) (GPIOState, error) {
	state, err := hif.GetGPIOValue(uint8(gpio))
	return GPIOState(state), err
}

func (w *WINC) isr(signal <-chan bool) {
	// Send the shutdown signal when this routine eventually returns
	defer func() { isrShutdownSignal <- true }()

	// Loop forever until the driver is reset
	for {
		select {
		case state := <-signal: // Wait for the signal from the interrupt
			if !state {
				// Stop this goroutine
				println("(", time.Now().String(), ") Shutting down ISR")
				return
			}

			// Wake the chip
			err := hif.ChipWake()

			if err == nil {
				// Handle the interrupt
				err = hif.Isr()

				if err != nil {
					println(err.Error())
				}

				// Sleep the chip
				if err = hif.ChipSleep(); err != nil {
					println(err.Error())
				}
			} else {
				println(err.Error())
			}
		}
	}
}

func IrqHandler(hal.Pin) {
	select {
	case isrSignal <- true:
	// Unblock the interrupt service (go)routine
	default:
		return
	}
}
