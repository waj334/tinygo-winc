package winc

import (
	"bytes"
	"encoding/binary"
	"github.com/waj334/tinygo-winc/debug"
	"github.com/waj334/tinygo-winc/protocol"
	"io"
)

type State uint8

const (
	stateWaitOpCode State = iota + 1
	stateWaitHeader
	stateProcessCommand
	stateWaitPayload
)

const (
	cmdBufferSize = 2048
	cmdHeaderSize = 12

	responseNack    byte = 0x5a
	responseIdVarBr byte = 0x5b
	responseAck     byte = 0xac
)

const (
	commandReadRegWithRet uint8 = iota
	commandWriteReg
	commandReadBlock
	commandWriteBlock
	commandReconfigure
)

type Serial interface {
	io.ReadWriter
	ReadByte() (byte, error)
	WriteByte(b byte) error
}

type SerialBridge struct {
	buf    *bytes.Buffer
	drv    *WINC
	serial Serial
	state  State

	cmdType uint8
	cmdSize uint16
	cmdAddr uint32
	cmdVal  uint32

	payloadLength uint16
}

func NewSerialBridge(serial Serial, drv *WINC) (*SerialBridge, error) {
	// Create the HIF
	drv.hif = protocol.CreateHif(drv.SPI, drv.CS)

	// Put WINC into download mode
	if err := drv.hif.InitDownload(); err != nil {
		return nil, err
	}

	// Allocate buffer for incoming commands
	cmdBuf := make([]byte, 0, cmdBufferSize)
	return &SerialBridge{
		buf:    bytes.NewBuffer(cmdBuf),
		drv:    drv,
		serial: serial,
		state:  stateWaitOpCode,
	}, nil
}

func (s *SerialBridge) Process() error {
	switch s.state {
	case stateWaitOpCode:
		// Receive the opcode
		var opcode byte
		var err error
		if opcode, err = s.serial.ReadByte(); err != nil {
			//println("serial.ReadByte:", err)
			break
		}

		debug.DEBUG("Opcode is %#x", opcode)

		switch opcode {
		case 0x12:
			s.serial.WriteByte(responseIdVarBr)

			// Reset buffer
			s.buf.Reset()

		case 0xa5:
			s.state = stateWaitHeader
		default:
			break
		}

	case stateWaitHeader:
		// Read header bytes
		s.buf.ReadFrom(s.serial)

		// Keep reading header bytes until exactly 12 bytes are read
		if s.buf.Len() != cmdHeaderSize {
			break
		}

		// Process the header
		if s.processHeader(s.buf.Bytes()) {
			// Ack the incoming command
			s.serial.WriteByte(responseAck)

			if s.payloadLength > 0 {
				// Go to payload wait state
				s.state = stateWaitPayload

				// Reset the buffer
				s.buf.Reset()
			} else {
				// Go to the command processing state
				s.state = stateProcessCommand
			}
		}
	case stateProcessCommand:
		// Process the command
		s.processCommand()

		// Go to the wait for incoming command state
		s.state = stateWaitOpCode
	case stateWaitPayload:
		// Receive the payload
		s.buf.ReadFrom(s.serial)

		// Begin processing the command once all the expected payload bytes are received
		if s.buf.Len() == int(s.payloadLength) {
			// Goto process command state
			s.state = stateProcessCommand
		}
	default:
		debug.DEBUG("Entered unknown state")
	}

	return nil
}

func (s *SerialBridge) processHeader(header []byte) bool {
	var checksum uint8

	if len(header) == 0 {
		return false
	}

	// Calculate checksum
	for i := range header {
		checksum ^= header[i]
	}

	// Checksum must be exactly zero
	if checksum != 0 {
		return false
	}

	s.cmdType = header[0]
	s.cmdSize = binary.BigEndian.Uint16(header[2:])
	s.cmdAddr = binary.BigEndian.Uint32(header[4:])
	s.cmdVal = binary.BigEndian.Uint32(header[8:])

	if s.cmdType == commandWriteBlock {
		s.payloadLength = s.cmdSize
		if s.payloadLength > cmdBufferSize {
			return false
		}
	} else {
		s.payloadLength = 0
	}

	return true
}

func (s *SerialBridge) processCommand() {
	switch s.cmdType {
	case commandReadRegWithRet:
		data := [4]byte{}

		// Read the value stored in the register at the current command's address
		regVal, _ := s.drv.hif.ReadRegister(s.cmdAddr)

		// Format the register value
		data[0] = byte((regVal >> 24) & 0xff)
		data[1] = byte((regVal >> 16) & 0xff)
		data[2] = byte((regVal >> 8) & 0xff)
		data[3] = byte((regVal) & 0xff)

		// Write the register value to the serial interface
		s.serial.Write(data[:])
	case commandWriteReg:
		// Write the value to the register specified by the command
		s.drv.hif.WriteRegister(s.cmdAddr, s.cmdVal)
	case commandReadBlock:
		// Allocate memory for block
		data := make([]byte, s.cmdSize)

		// Read the data block at the address
		if err := s.drv.hif.ReadBlock(s.cmdAddr, data); err != nil {
			return
		}

		// Write the data block to serial interface
		s.serial.Write(data)
	case commandWriteBlock:
		// Write the buffered data to the address
		if err := s.drv.hif.WriteBlock(s.cmdAddr, s.buf.Bytes()); err != nil {
			// NACK the command
			s.serial.WriteByte(responseNack)
		}

		// ACK the command
		s.serial.WriteByte(responseAck)

		// Echo the written data on the serial interface
		s.serial.Write(s.buf.Bytes())
	case commandReconfigure:
		// TODO: The original implementation changes the baud rate of the the serial interface. Implement this later if
		//       this is actually needed.
	}
}
