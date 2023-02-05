package protocol

import (
	"bytes"
	"encoding/binary"
)

type hifHeader struct {
	groupId byte
	opcode  byte
	length  uint16
}

func (h *hifHeader) bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 4))
	buf.WriteByte(h.groupId)
	buf.WriteByte(h.opcode)
	binary.Write(buf, binary.LittleEndian, h.length)
	return buf.Bytes()
}
