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
	"machine"
	"runtime/volatile"
	"sync"
	"time"

	"tinygo.org/x/drivers"

	"github.com/waj334/tinygo-winc/debug"
)

type (
	GroupId  uint8
	OpcodeId uint8
)

const (
	GroupMax     = 9
	HifHdrOffset = 8

	hifMaxPacketSize = 1600 - 4
)

const (
	OpcodeReqConfigPkt OpcodeId = 0
	OpcodeReqDataPkt   OpcodeId = 0x80 /*BIT7*/
)

var (
	rxSize    uint32
	rxAddress uint32
	rxDone    uint8
	chipId    uint32
	callbacks [GroupMax]IsrCallback
)

type IsrCallback func(id OpcodeId, sz uint16, address uint32) (any, error)

type Event struct {
	Group  GroupId
	Opcode OpcodeId
	Data   any
}

type Hif struct {
	t             transport
	eventChannels []chan Event
	mutex         sync.Mutex
}

func CreateHif(spi drivers.SPI, cs machine.Pin) Hif {
	return Hif{
		t: transport{
			spi: spi,
			cs:  cs,
		},
	}
}

func (hif *Hif) Init() (err error) {
	callbacks = [GroupMax]IsrCallback{}

	if err = hif.t.init(); err != nil {
		return err
	}

	if chipId, err := hif.GetChipId(); err != nil {
		return err
	} else if (chipId >> 16) != 0x15 { // TODO: Allow WINC3400
		return errIncompatibleVersion
	}

	if err = hif.waitForBootrom(); err != nil {
		return errBootromFailed
	}

	if err = hif.waitForFirmwareStart(); err != nil {
		return errFirmwareFailed
	}

	if err = hif.enableInterrupts(); err != nil {
		return errInterruptsEnableFailed
	}

	return nil
}

func (hif *Hif) InitDownload() (err error) {
	if err = hif.t.init(); err != nil {
		return err
	}

	chipId, err := hif.GetChipId()
	if err != nil {
		return
	}

	if (chipId & 0xffff0000) != 0x300000 {
		if err = hif.Halt(); err != nil {
			return
		}
	}

	// Init the transport
	if err = hif.t.init(); err != nil {
		return
	}

	// Disable all interrupts
	if err = hif.t.WriteRegister(0x20300, 0); err != nil {
		return
	}

	return
}

func (hif *Hif) Shutdown() {
	// Close all event channels
	for _, e := range hif.eventChannels {
		select {
		case <-e: // Drain the channel
			close(e)
		default: // The channel is already drained
			close(e)
		}
	}

	rxSize = 0
	rxAddress = 0
	rxDone = 0
	chipId = 0
}

func (hif *Hif) RegisterCallback(group GroupId, callback IsrCallback) {
	hif.mutex.Lock()
	defer hif.mutex.Unlock()

	callbacks[group] = callback
}

func (hif *Hif) OpenEventChannel() <-chan Event {
	hif.mutex.Lock()
	defer hif.mutex.Unlock()

	c := make(chan Event, 1)
	hif.eventChannels = append(hif.eventChannels, c)
	return c
}

func (hif *Hif) ChipWake() (err error) {
	hif.mutex.Lock()
	defer hif.mutex.Unlock()
	return hif.chipWakeInternal()
}

func (hif *Hif) chipWakeInternal() (err error) {
	if volatile.LoadUint8(&rxDone) != 0 {
		// Chip already wake
		return nil
	}

	if err = hif.t.WriteRegister(_HOST_CORT_COMM, _NBIT0); err != nil {
		return err
	}
	if err = hif.t.WriteRegister(_WAKE_CLOCK_REG, _NBIT1); err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 3)

	// Receive clock enabled register until bit 2 is 1
	for retries := 0; retries < 10; retries++ {
		var reg uint32
		if reg, err = hif.t.ReadRegister(_CLOCKS_EN_REG); err != nil {
			return err
		} else if reg&_NBIT2 != 0 {
			// Reset the bus
			//hif.t.reset()

			return nil
		}

		// Sleep a while before trying again
		time.Sleep(time.Millisecond * 2)
	}

	return errChipWakeFail
}

func (hif *Hif) ChipSleep() error {
	hif.mutex.Lock()
	defer hif.mutex.Unlock()
	return hif.chipSleepInternal()
}

func (hif *Hif) chipSleepInternal() error {
	for {
		result, err := hif.t.ReadRegister(_CORT_HOST_COMM)
		if err != nil {
			return err
		} else if result&_NBIT0 == 0 {
			break
		}
	}

	// Clear bit 1
	result, err := hif.t.ReadRegister(_WAKE_CLOCK_REG)
	if err != nil {
		return err
	}

	if result&_NBIT1 != 0 {
		result &= ^uint32(_NBIT1)
		if err = hif.t.WriteRegister(_WAKE_CLOCK_REG, result); err != nil {
			return err
		}
	}

	result, err = hif.t.ReadRegister(_HOST_CORT_COMM)
	if err != nil {
		return err
	}

	if result&_NBIT0 != 0 {
		result &= ^uint32(_NBIT0)
		err = hif.t.WriteRegister(_HOST_CORT_COMM, result)
		if err != nil {
			return err
		}
	}

	return nil
}

func (hif *Hif) GetChipId() (uint32, error) {
	hif.mutex.Lock()
	defer hif.mutex.Unlock()

	if chipId == 0 {
		var err error
		// Receive the chip ID
		chipId, err = hif.t.ReadRegister(_NMI_CHIPID)
		if err != nil {
			return 0, err
		}

		var rfrevid uint32
		if rfrevid, err = hif.t.ReadRegister(0x13F4); err != nil {
			return 0, err
		}

		if chipId == 0x1002A0 {
			if rfrevid != 0x1 {
				chipId = 0x1002A1
			}
		} else if chipId == 0x1002B0 {
			if rfrevid == 4 {
				chipId = 0x1002B1
			} else {
				chipId = 0x1002B2
			}
		} else if chipId == 0x1000F0 {
			if chipId, err = hif.t.ReadRegister(0x3B0000); err != nil {
				return 0, err
			}
		}

		chipId &= ^uint32(0x0F0000)
		chipId |= 0x050000
	}

	return chipId, nil
}

func (hif *Hif) Receive(address uint32, data []byte, done bool) (err error) {
	hif.mutex.Lock()
	defer hif.mutex.Unlock()

	if address == 0 || data == nil || len(data) == 0 {
		if done {
			return hif.setRxDone()
		} else {
			return errOperationFailed
		}
	}

	// NOTE: This is set by the ISR
	rxSize := volatile.LoadUint32(&rxSize)
	rxAddress := volatile.LoadUint32(&rxAddress)
	length := uint32(len(data))

	if length > rxSize {
		return errBadMemoryAlloc
	}

	if address < rxAddress || (address+length) > (rxAddress+rxSize) {
		return errMessageTooLong
	}

	// Receive the packet
	if err = hif.t.ReadBlock(address, data); err != nil {
		return err
	}

	// Is this the last packet?
	if (rxAddress+rxSize)-(address+length) <= 0 || done {
		// Set RX done
		return hif.setRxDone()
	}

	return nil
}

func (hif *Hif) Send(group GroupId, opcode OpcodeId, control, data []byte, offset uint16) (err error) {
	hif.mutex.Lock()
	defer hif.mutex.Unlock()

	debug.DEBUG("HIF: Send - BEGIN")
	defer debug.DEBUG("HIF: Send - END")

	strHif := hifHeader{
		groupId: byte(group),
		opcode:  byte(opcode) & ^byte(_NBIT7),
		length:  HifHdrOffset,
	}

	if dataLen := len(data); dataLen > 0 {
		// Add length of data and the offset
		strHif.length += offset + uint16(dataLen)
	} else {
		// Add just the length of the control struct
		strHif.length += uint16(len(control))
	}

	if strHif.length < hifMaxPacketSize {
		// Wake the chip before sending the HIF packet
		if err = hif.chipWakeInternal(); err != nil {
			return err
		}

		var reg uint32
		reg |= uint32(group)
		reg |= uint32(opcode) << 8
		reg |= uint32(strHif.length) << 16

		if err = hif.WriteRegister(_NMI_STATE_REG, reg); err != nil {
			return err
		}

		if err = hif.WriteRegister(_WIFI_HOST_RCV_CTRL_2, _NBIT1); err != nil {
			return err
		}

		var dmaAddr uint32
		for cnt := 0; cnt < 1000; cnt++ {
			if reg, err = hif.ReadRegister(_WIFI_HOST_RCV_CTRL_2); err != nil {
				return err
			}

			// Throttle requests after 500 retries
			if cnt > 500 {
				if cnt < 501 {
					debug.DEBUG("HIF: Slowing down...")
				}
				time.Sleep(time.Nanosecond)
			}

			if reg&_NBIT1 == 0 {
				if dmaAddr, err = hif.ReadRegister(_WIFI_HOST_RCV_CTRL_4); err != nil {
					return err
				}
				break
			}
		}

		if dmaAddr != 0 {
			currentAddr := dmaAddr
			debug.DEBUG("HIF: Writing HIF header.")
			debug.DEBUG("HIF: strHif.U8Gid: %#v", strHif.groupId)
			debug.DEBUG("HIF: strHif.U8Opcode: %v", strHif.opcode)
			debug.DEBUG("HIF: strHif.U16Length: %v", strHif.length)
			debug.DEBUG("HIF: currentAddr: %#x", currentAddr)

			// NOTE: The HIF header is transmitted in 8 bytes even though the struct itself is 4 bytes.
			var headerBytes [HifHdrOffset]byte
			copy(headerBytes[:], strHif.bytes())

			if err = hif.t.WriteBlock(currentAddr, headerBytes[:]); err != nil {
				return err
			}

			debug.DEBUG("HIF: Done writing HIF header.")
			currentAddr += HifHdrOffset

			if len(control) > 0 {
				debug.DEBUG("HIF: Writing control bytes.")
				debug.DEBUG("HIF: length: %v", len(control))
				debug.DEBUG("HIF: currentAddr: %#x", currentAddr)

				if err = hif.WriteBlock(currentAddr, control); err != nil {
					return err
				}

				currentAddr += uint32(len(control))
				debug.DEBUG("HIF: Done writing control bytes")
			}

			if len(data) > 0 {
				debug.DEBUG("HIF: Writing data bytes.")
				debug.DEBUG("HIF: length: %v", len(data))
				debug.DEBUG("HIF: currentAddr: %#x", currentAddr)

				// Advance address before writing the data so that it is written at the correct offset.
				currentAddr += uint32(offset - uint16(len(control)))
				if err = hif.WriteBlock(currentAddr, data); err != nil {
					return err
				}

				debug.DEBUG("HIF: Done writing data bytes.")
			}

			reg = dmaAddr << 2
			reg |= _NBIT1
			if err = hif.WriteRegister(_WIFI_HOST_RCV_CTRL_3, reg); err != nil {
				return err
			}
		} else {
			debug.DEBUG("HIF: Failed to receive DMA address")
			if err = hif.ChipSleep(); err != nil {
				debug.DEBUG("HIF: hif.ChipSleep - %v", err)
				return err
			}
			return errBadMemoryAlloc
		}
	} else {
		return errMessageTooLong
	}

	hif.chipSleepInternal()
	return
}

func (hif *Hif) waitForBootrom() error {
	// Wait for efuse loading
	for {
		if value, _ := hif.t.ReadRegister(0x1014); value&0x80000000 != 0 {
			break
		}
		time.Sleep(time.Millisecond)
	}

	value, _ := hif.t.ReadRegister(_WAIT_FOR_HOST_REG)
	value &= 0x1

	timeout := time.Now().Add(time.Millisecond * 0x2000)

	// Check whether waiting on the host should be skipped
	if value == 0 {
		for value, _ = hif.t.ReadRegister(_BOOTROM_REG); uint64(value) != _FINISH_BOOT_ROM; {
			if time.Now().After(timeout) {
				return errFirmwareLoadFailed
			}

			time.Sleep(time.Millisecond)
		}
	}

	// Write the version info
	err := hif.t.WriteRegister(_NMI_STATE_REG, _VERSION)
	if err != nil {
		return err
	}

	if chipId, _ := hif.GetChipId(); chipId > _REV_3A0 {
		hif.applyChipConfig(_HAVE_USE_PMU_BIT)
	} else {
		hif.applyChipConfig(0)
	}

	// Start the firmware
	return hif.t.WriteRegister(_BOOTROM_REG, uint32(_START_FIRMWARE))
}

func (hif *Hif) waitForFirmwareStart() error {
	timeout := time.Now().Add(time.Millisecond * 0x2000)
	for {
		if value, _ := hif.t.ReadRegister(_NMI_STATE_REG); uint64(value) != _FINISH_INIT_STATE {
			if time.Now().After(timeout) {
				return errFirmwareTimeout
			}
		} else {
			break
		}

		time.Sleep(time.Millisecond * 1000)
	}

	// Clear the state register
	hif.t.WriteRegister(_NMI_STATE_REG, 0)

	return nil
}

func (hif *Hif) enableInterrupts() error {
	// Interrupt pin mux select
	value, err := hif.t.ReadRegister(_NMI_PIN_MUX_0)
	if err != nil {
		return err
	}

	value |= 1 << 8
	if err = hif.t.WriteRegister(_NMI_PIN_MUX_0, value); err != nil {
		return err
	}

	// Enable the interrupt for the pin
	if value, err = hif.t.ReadRegister(_NMI_INTR_ENABLE); err != nil {
		return err
	}

	value |= 1 << 16
	if err = hif.t.WriteRegister(_NMI_INTR_ENABLE, value); err != nil {
		return err
	}

	return nil
}

func (hif *Hif) applyChipConfig(conf uint32) error {
	conf |= _HAVE_RESERVED1_BIT
	for {
		hif.t.WriteRegister(_NMI_GP_REG_1, conf)
		if conf != 0 {
			if value, err := hif.t.ReadRegister(_NMI_GP_REG_1); err == nil && value == conf {
				break
			}
		} else {
			break
		}
	}

	// The original driver had no fail state
	return nil
}

func (hif *Hif) writeHeader(groupId, opcode uint8, length uint16, address uint32) (err error) {
	data := [4]byte{
		groupId,
		opcode & (^uint8(_NBIT7)),
		byte(length >> 8), // TODO: Account for endianess when encoding the length
		byte(length),
	}

	// Write the header
	if err = hif.t.WriteBlock(address, data[:]); err != nil {
		return
	}

	return
}

func (hif *Hif) readHeader(address uint32) (group GroupId, opcode OpcodeId, length uint16, err error) {
	var data [4]byte
	if err = hif.t.ReadBlock(address, data[:]); err != nil {
		return
	}

	group = GroupId(data[0])
	opcode = OpcodeId(data[1])

	// TODO: Account for endianess when decoding the length
	length = uint16(data[2]) | (uint16(data[3]) << 8)

	return
}

func (hif *Hif) setRxDone() (err error) {
	volatile.StoreUint8(&rxDone, 0)
	var reg uint32
	if reg, err = hif.t.ReadRegister(_WIFI_HOST_RCV_CTRL_0); err != nil {
		return err
	}

	// Set RX done
	reg |= _NBIT1
	if err = hif.t.WriteRegister(_WIFI_HOST_RCV_CTRL_0, reg); err != nil {
		return err
	}

	return
}

func (hif *Hif) Isr() (err error) {
	// Lock the mutex to prevent other goroutines from sending frames will this ISR is processing
	var once sync.Once
	hif.mutex.Lock()
	defer once.Do(hif.mutex.Unlock)

	var size uint16

	// Receive RX interrupt state
	var reg uint32
	if reg, err = hif.t.ReadRegister(_WIFI_HOST_RCV_CTRL_0); err != nil {
		return
	}

	// Has the RX interrupt been received
	if reg&0x1 != 0 {
		// Clear RX interrupt
		reg &= ^uint32(_NBIT0)
		if err = hif.t.WriteRegister(_WIFI_HOST_RCV_CTRL_0, reg); err != nil {
			return
		}
	}

	// Set the RX done state
	volatile.StoreUint8(&rxDone, 1)

	// Set the size
	size = uint16(reg>>2) & 0xFFF
	if size > 0 {
		// Start the bus transfer
		var address uint32
		if address, err = hif.t.ReadRegister(_WIFI_HOST_RCV_CTRL_1); err != nil {
			return
		}

		volatile.StoreUint32(&rxAddress, address)
		volatile.StoreUint32(&rxSize, uint32(size))

		// Receive the header
		var group GroupId
		var opcode OpcodeId
		var length uint16
		if group, opcode, length, err = hif.readHeader(address); err != nil {
			return
		}

		if size-length > 4 {
			// The packet is likely corrupted
			return errOperationFailed
		}

		// Unlock the mutex so other goroutines can run/continue
		once.Do(hif.mutex.Unlock)

		// Execute the respective callback functions based on the header
		var data any
		var callbackErr error
		if fn := callbacks[group]; fn != nil {
			if data, callbackErr = fn(opcode, length-8, address+8); callbackErr == nil && data != nil {
				// Emit event
				e := Event{
					Group:  group,
					Opcode: opcode,
					Data:   data,
				}

				for i := 0; i < len(hif.eventChannels); i++ {
					eventChan := hif.eventChannels[i]
					select {
					case eventChan <- e:
						// Proceed
					default:
						// Attempt to drain the channel
						if _, ok := <-eventChan; !ok {
							// This channel is closed. Remove from list
							hif.eventChannels = append(hif.eventChannels[:i], hif.eventChannels[i+1:]...)
						} else {
							// Signal the event
							eventChan <- e
						}
					}
				}
			}
		}

		volatile.LoadUint8(&rxDone)
		if rxDone != 0 {
			if err = hif.setRxDone(); err != nil {
				return err
			}
		}

		return callbackErr
	}

	return
}

func (hif *Hif) SetGPIODirection(gpio, direction uint8) (err error) {
	hif.mutex.Lock()
	defer hif.mutex.Unlock()

	var value uint32
	if value, err = hif.t.ReadRegister(0x20108); err != nil {
		return
	}

	if direction != 0 {
		value |= 1 << gpio
	} else {
		value &= ^(1 << gpio)
	}

	return hif.t.WriteRegister(0x20108, value)
}

func (hif *Hif) SetGPIOValue(gpio, state uint8) (err error) {
	hif.mutex.Lock()
	defer hif.mutex.Unlock()

	var value uint32
	if value, err = hif.t.ReadRegister(0x20100); err != nil {
		return
	}

	if state != 0 {
		value |= 1 << gpio
	} else {
		value &= ^(1 << gpio)
	}

	return hif.t.WriteRegister(0x20100, value)
}

func (hif *Hif) GetGPIOValue(gpio uint8) (state uint8, err error) {
	hif.mutex.Lock()
	defer hif.mutex.Unlock()

	var value uint32
	if value, err = hif.t.ReadRegister(0x142C); err != nil {
		return
	}

	state = uint8(value>>gpio) & 0x01

	return
}

func (hif *Hif) Halt() error {
	reg, err := hif.t.ReadRegister(0x1118)
	if err != nil {
		return err
	}

	reg |= 1 << 0

	if err = hif.t.WriteRegister(0x1118, reg); err != nil {
		return err
	}

	if reg, err = hif.t.ReadRegister(_NMI_GLB_RESET_0); err != nil {
		return err
	}

	if (reg & (1 << 10)) == (1 << 10) {
		reg &= ^uint32(1 << 10)
		if err = hif.t.WriteRegister(_NMI_GLB_RESET_0, reg); err != nil {
			return err
		}

		if reg, err = hif.t.ReadRegister(_NMI_GLB_RESET_0); err != nil {
			return err
		}
	}

	return nil
}

func (hif *Hif) ReadRegister(address uint32) (uint32, error) {
	return hif.t.ReadRegister(address)
}

func (hif *Hif) WriteRegister(address, value uint32) error {
	return hif.t.WriteRegister(address, value)
}

func (hif *Hif) ReadBlock(address uint32, data []byte) (err error) {
	return hif.t.ReadBlock(address, data)
}

func (hif *Hif) WriteBlock(address uint32, data []byte) (err error) {
	return hif.t.WriteBlock(address, data)
}
