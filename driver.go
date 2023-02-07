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
	"net"
	"sync"
	"time"

	"machine"
	"tinygo.org/x/drivers"

	"github.com/waj334/tinygo-winc/debug"
	"github.com/waj334/tinygo-winc/protocol"
)

type WINC struct {
	SPI       drivers.SPI
	CS        machine.Pin
	IRQ       machine.Pin
	EnablePin machine.Pin
	ResetPin  machine.Pin

	EccProvider EccProvider

	wifiState WifiState
	ipAddr    net.IPNet

	hif protocol.Hif

	initialized       bool
	isrSignal         chan bool
	isrShutdownSignal chan bool

	callbackChan chan any

	sockets             [maxSocket]*Socket
	sessionCounterMutex sync.Mutex
	sessionCounter      uint16
	SocketBufferLength  int

	mutex sync.Mutex
}

func (w *WINC) Initialize() (err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.initialized {
		// Create the hardware interface abstraction layer
		w.hif = protocol.CreateHif(w.SPI, w.CS)

		// Initialize the HAL
		err = w.hif.Init()
		if err != nil {
			return
		}

		// Register interrupt callbacks
		w.hif.RegisterCallback(GroupWIFI, w.wifiCallback)
		w.hif.RegisterCallback(GroupIP, w.socketCallback)
		w.hif.RegisterCallback(GroupSSL, w.sslCallback)

		// Start the interrupt service (go)routine
		w.isrSignal = make(chan bool, 1)
		w.isrShutdownSignal = make(chan bool, 1)
		go w.isr(w.isrSignal)

		// Enable the interrupt
		w.setInterruptEnabled(true)

		// Create the callback channel
		w.callbackChan = make(chan any, 1)

		// set up sockets
		w.sessionCounter = 1
		if w.SocketBufferLength == 0 {
			// Set to default
			w.SocketBufferLength = 2048
		}

		w.initialized = true
	}

	return
}

func (w *WINC) OpenEventChannel() <-chan protocol.Event {
	return w.hif.OpenEventChannel()
}

// Reset the SoC by driving chip enable and reset pins low then high
func (w *WINC) Reset() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// Unset the interrupt handler
	w.setInterruptEnabled(false)

	// Stop the ISR routine
	if w.isrSignal != nil {
		select {
		case w.isrSignal <- false:
			close(w.isrSignal)
		default:
			close(w.isrSignal)
		}

		// Wait for isr to be fully shutdown
		<-w.isrShutdownSignal
		close(w.isrShutdownSignal)

		w.isrSignal = nil
	}

	// Shutdown driver.hif
	w.hif.Shutdown()

	// Reset sockets
	w.sockets = [maxSocket]*Socket{}
	//currentSocket = nil

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

	w.initialized = false
}

func (w *WINC) SetGPIODirection(gpio GPIOType, direction GPIODirection) error {
	return w.hif.SetGPIODirection(uint8(gpio), uint8(direction))
}

func (w *WINC) SetGPIOState(gpio GPIOType, state GPIOState) error {
	return w.hif.SetGPIOValue(uint8(gpio), uint8(state))
}

func (w *WINC) GetGPIOState(gpio GPIOType) (GPIOState, error) {
	state, err := w.hif.GetGPIOValue(uint8(gpio))
	return GPIOState(state), err
}

func (w *WINC) setInterruptEnabled(on bool) {
	if on {
		w.IRQ.SetInterrupt(machine.PinFalling, w.irqHandler)
	} else {
		w.IRQ.SetInterrupt(machine.PinFalling, nil)
	}
}

func (w *WINC) isr(signal <-chan bool) {
	// Send the shutdown signal when this routine eventually returns
	defer func() { w.isrShutdownSignal <- true }()

	// Loop forever until the driver is reset
	for {
		select {
		case state := <-signal: // Wait for the signal from the interrupt
			if !state {
				// Stop this goroutine
				return
			}

			// Wake the chip
			err := w.hif.ChipWake()

			if err == nil {
				// Handle the interrupt
				err = w.hif.Isr()
				if err != nil {
					debug.DEBUG("ISR error: %v", err)
				}
				// Sleep the chip
				if err = w.hif.ChipSleep(); err != nil {
					debug.DEBUG("ISR error: %v", err)
				}
			} else {
				debug.DEBUG("ISR error: %v", err)
			}
		}
	}
}

func (w *WINC) irqHandler(machine.Pin) {
	select {
	case w.isrSignal <- true:
	// Unblock the interrupt service (go)routine
	default:
		return
	}
}
